module "vpc" {
  source = "terraform-aws-modules/vpc/aws"

  name = var.vpc.name
  cidr = var.vpc.cidr

  azs                                   = var.vpc.azs
  private_subnets                       = var.vpc.private_subnets
  public_subnets                        = var.vpc.public_subnets
  database_subnets                      = var.vpc.database_subnets
  elasticache_subnets                   = var.vpc.elasticache_subnets
  create_database_subnet_group          = var.vpc.create_database_subnet_group
  create_database_subnet_route_table    = var.vpc.create_database_subnet_route_table
  create_elasticache_subnet_group       = var.vpc.create_elasticache_subnet_group
  create_elasticache_subnet_route_table = var.vpc.create_elasticache_subnet_route_table
  enable_vpn_gateway                    = var.vpc.enable_vpn_gateway
  one_nat_gateway_per_az                = var.vpc.one_nat_gateway_per_az
  single_nat_gateway                    = var.vpc.single_nat_gateway
  enable_nat_gateway                    = var.vpc.enable_nat_gateway
}

module "byo-vpc" {
  source = "./byo-vpc"
  vpc_id = module.vpc.vpc_id
}
