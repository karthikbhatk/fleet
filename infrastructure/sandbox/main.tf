terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.10.0"
    }
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 2.16.0"
    }
    git = {
      source  = "paultyng/git"
      version = "~> 0.1.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.1.2"
    }
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 3.18.0"
    }
  }
  backend "s3" {}
}

provider "aws" {
  region = "us-east-2"
  default_tags {
    tags = {
      environment = "fleet-demo-${terraform.workspace}"
      terraform   = "https://github.com/fleetdm/fleet/tree/main/infrastructure/demo"
      state       = "s3://fleet-loadtesting-tfstate/demo-environment"
    }
  }
}
provider "aws" {
  alias  = "replica"
  region = "us-west-1"
  default_tags {
    tags = {
      environment = "fleet-demo-${terraform.workspace}"
      terraform   = "https://github.com/fleetdm/fleet/tree/main/infrastructure/demo"
      state       = "s3://fleet-loadtesting-tfstate/demo-environment"
    }
  }
}

provider "cloudflare" {}

provider "random" {}

data "aws_ecr_authorization_token" "token" {}
provider "docker" {
  # Configuration options
  registry_auth {
    address  = "${data.aws_caller_identity.current.account_id}.dkr.ecr.us-east-2.amazonaws.com"
    username = data.aws_ecr_authorization_token.token.user_name
    password = data.aws_ecr_authorization_token.token.password
  }
}

provider "git" {}

data "aws_caller_identity" "current" {}

data "git_repository" "tf" {
  path = "${path.module}/../../"
}

locals {
  prefix      = "sandbox-prod"
  base_domain = "sandbox.fleetdm.com"
}

data "aws_iam_policy_document" "kms" {
  statement {
    actions = ["kms:*"]
    principals {
      type        = "AWS"
      identifiers = ["arn:aws:iam::${data.aws_caller_identity.current.account_id}:root"]
    }
    resources = ["*"]
  }
  statement {
    actions = [
      "kms:Encrypt*",
      "kms:Decrypt*",
      "kms:ReEncrypt*",
      "kms:GenerateDataKey*",
      "kms:Describe*",
    ]
    resources = ["*"]
    principals {
      type = "Service"
      # TODO hard coded region
      identifiers = ["logs.us-east-2.amazonaws.com"]
    }
  }
}

resource "aws_kms_key" "main" {
  policy              = data.aws_iam_policy_document.kms.json
  enable_key_rotation = true
}

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "3.12.0"

  name = local.prefix
  cidr = "10.11.0.0/16"

  # TODO hard coded AZs
  azs                 = ["us-east-2a", "us-east-2b", "us-east-2c"]
  private_subnets     = ["10.11.16.0/20", "10.11.32.0/20", "10.11.48.0/20"]
  public_subnets      = ["10.11.128.0/24", "10.11.129.0/24", "10.11.130.0/24"]
  database_subnets    = ["10.11.131.0/24", "10.11.132.0/24", "10.11.133.0/24"]
  elasticache_subnets = ["10.11.134.0/24", "10.11.135.0/24", "10.11.136.0/24"]

  create_database_subnet_group       = false
  create_database_subnet_route_table = true

  create_elasticache_subnet_group       = true
  create_elasticache_subnet_route_table = true

  enable_vpn_gateway     = false
  one_nat_gateway_per_az = false

  single_nat_gateway = true
  enable_nat_gateway = true
}

module "shared-infrastructure" {
  source                  = "./SharedInfrastructure"
  prefix                  = local.prefix
  vpc                     = module.vpc
  allowed_security_groups = [module.pre-provisioner.lambda_security_group.id]
  eks_allowed_roles       = [module.pre-provisioner.lambda_role, module.jit-provisioner.deprovisioner_role]
  base_domain             = local.base_domain
  kms_key                 = aws_kms_key.main
}

module "pre-provisioner" {
  source            = "./PreProvisioner"
  prefix            = local.prefix
  vpc               = module.vpc
  kms_key           = aws_kms_key.main
  dynamodb_table    = aws_dynamodb_table.lifecycle-table
  remote_state      = module.remote_state
  mysql_secret      = module.shared-infrastructure.mysql_secret
  eks_cluster       = module.shared-infrastructure.eks_cluster
  redis_cluster     = module.shared-infrastructure.redis_cluster
  ecs_cluster       = aws_ecs_cluster.main
  base_domain       = local.base_domain
  installer_bucket  = module.shared-infrastructure.installer_bucket
  oidc_provider_arn = module.shared-infrastructure.oidc_provider_arn
  oidc_provider     = module.shared-infrastructure.oidc_provider
}

module "jit-provisioner" {
  source         = "./JITProvisioner"
  prefix         = local.prefix
  vpc            = module.vpc
  kms_key        = aws_kms_key.main
  dynamodb_table = aws_dynamodb_table.lifecycle-table
  remote_state   = module.remote_state
  mysql_secret   = module.shared-infrastructure.mysql_secret
  eks_cluster    = module.shared-infrastructure.eks_cluster
  redis_cluster  = module.shared-infrastructure.redis_cluster
  alb_listener   = module.shared-infrastructure.alb_listener
  ecs_cluster    = aws_ecs_cluster.main
  base_domain    = local.base_domain
}

module "monitoring" {
  source         = "./Monitoring"
  prefix         = local.prefix
  slack_webhook  = var.slack_webhook
  kms_key        = aws_kms_key.main
  lb             = module.shared-infrastructure.lb
  jitprovisioner = module.jit-provisioner.jitprovisioner
  deprovisioner  = module.jit-provisioner.deprovisioner
  dynamodb_table = aws_dynamodb_table.lifecycle-table
}

resource "aws_dynamodb_table" "lifecycle-table" {
  name         = "${local.prefix}-lifecycle"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "ID"

  server_side_encryption {
    enabled     = true
    kms_key_arn = aws_kms_key.main.arn
  }
  point_in_time_recovery {
    enabled = true
  }

  attribute {
    name = "ID"
    type = "S"
  }

  attribute {
    name = "State"
    type = "S"
  }

  attribute {
    name = "redis_db"
    type = "N"
  }

  global_secondary_index {
    name            = "RedisDatabases"
    hash_key        = "redis_db"
    projection_type = "KEYS_ONLY"
  }
  global_secondary_index {
    name            = "FleetState"
    hash_key        = "State"
    projection_type = "ALL"
  }
}

module "remote_state" {
  source = "nozaq/remote-state-s3-backend/aws"
  tags   = {}

  providers = {
    aws         = aws
    aws.replica = aws.replica
  }
}

resource "aws_ecs_cluster" "main" {
  name = local.prefix

  setting {
    name  = "containerInsights"
    value = "enabled"
  }
}

variable "slack_webhook" {}
