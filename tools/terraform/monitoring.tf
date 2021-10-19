// sns topic to send cloudwatch alarms to
resource "aws_sns_topic" "cloudwatch_alarm_topic" {
  name = "cloudwatch-alarm-${terraform.workspace}"
}

resource "aws_sns_topic_policy" "default" {
  arn    = aws_sns_topic.cloudwatch_alarm_topic.arn
  policy = data.aws_iam_policy_document.sns_topic_policy.json
}

data "aws_iam_policy_document" "sns_topic_policy" {
  statement {
    sid = "AllowManageSNS"

    actions = [
      "sns:Subscribe",
      "sns:SetTopicAttributes",
      "sns:RemovePermission",
      "sns:Receive",
      "sns:Publish",
      "sns:ListSubscriptionsByTopic",
      "sns:GetTopicAttributes",
      "sns:DeleteTopic",
      "sns:AddPermission",
    ]

    effect    = "Allow"
    resources = [aws_sns_topic.cloudwatch_alarm_topic.arn]

    principals {
      type        = "AWS"
      identifiers = ["*"]
    }

    condition {
      test     = "StringEquals"
      variable = "AWS:SourceOwner"

      values = [data.aws_caller_identity.current.account_id]

    }
  }

  statement {
    sid       = "Allow CloudwatchEvents"
    actions   = ["sns:Publish"]
    resources = [aws_sns_topic.cloudwatch_alarm_topic.arn]

    principals {
      type        = "Service"
      identifiers = ["events.amazonaws.com"]
    }
  }

  statement {
    sid       = "Allow RDS Event Notification"
    actions   = ["sns:Publish"]
    resources = [aws_sns_topic.cloudwatch_alarm_topic.arn]

    principals {
      type        = "Service"
      identifiers = ["rds.amazonaws.com"]
    }
  }
}


// Database alarms
resource "aws_cloudwatch_metric_alarm" "cpu_utilization_too_high" {
  for_each            = toset(module.aurora_mysql.rds_cluster_instance_ids)
  alarm_name          = "rds_cpu_utilization_too_high-${each.key}-${terraform.workspace}"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = "1"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/RDS"
  period              = "300"
  statistic           = "Average"
  threshold           = 80
  alarm_description   = "Average database CPU utilization over last 5 minutes too high"
  alarm_actions       = [aws_sns_topic.cloudwatch_alarm_topic.arn]
  ok_actions          = [aws_sns_topic.cloudwatch_alarm_topic.arn]
  dimensions = {
    DBInstanceIdentifier = each.key
  }
}

resource "aws_db_event_subscription" "default" {
  name      = "rds-event-sub-${terraform.workspace}"
  sns_topic = aws_sns_topic.cloudwatch_alarm_topic.arn

  source_type = "db-instance"
  source_ids  = module.aurora_mysql.rds_cluster_instance_ids

  event_categories = [
    "failover",
    "failure",
    "low storage",
    "maintenance",
    "notification",
    "recovery",
  ]

  depends_on = [
    aws_sns_topic_policy.default
  ]
}

// ECS Alarms
resource "aws_cloudwatch_metric_alarm" "alb_healthyhosts" {
  alarm_name          = "backend-healthyhosts-${terraform.workspace}"
  comparison_operator = "LessThanThreshold"
  evaluation_periods  = "1"
  metric_name         = "HealthyHostCount"
  namespace           = "AWS/ApplicationELB"
  period              = "60"
  statistic           = "Minimum"
  threshold           = var.fleet_min_capacity
  alarm_description   = "This alarm indicates the number of Healthy Fleet hosts is lower than expected. Please investigate the load balancer \"${aws_alb.main.name}\" or the target group \"${aws_alb_target_group.main.name}\" and the fleet backend service \"${aws_ecs_service.fleet.name}\""
  actions_enabled     = "true"
  alarm_actions       = [aws_sns_topic.cloudwatch_alarm_topic.arn]
  ok_actions          = [aws_sns_topic.cloudwatch_alarm_topic.arn]
  dimensions = {
    TargetGroup  = aws_alb_target_group.main.arn_suffix
    LoadBalancer = aws_alb.main.arn_suffix
  }
}

// alarm for target response time (anomaly detection)
resource "aws_cloudwatch_metric_alarm" "target_response_time" {
  alarm_name                = "backend-target-response-time-${terraform.workspace}"
  comparison_operator       = "GreaterThanUpperThreshold"
  evaluation_periods        = "2"
  threshold_metric_id       = "e1"
  alarm_description         = "This alarm indicates the Fleet server response time is greater than it usually is. Please investigate the ecs service \"${aws_ecs_service.fleet.name}\" because the backend might need to be scaled up."
  alarm_actions             = [aws_sns_topic.cloudwatch_alarm_topic.arn]
  ok_actions                = [aws_sns_topic.cloudwatch_alarm_topic.arn]
  insufficient_data_actions = []

  metric_query {
    id          = "e1"
    expression  = "ANOMALY_DETECTION_BAND(m1)"
    label       = "TargetResponseTime (Expected)"
    return_data = "true"
  }

  metric_query {
    id          = "m1"
    return_data = "true"
    metric {
      metric_name = "TargetResponseTime"
      namespace   = "AWS/ApplicationELB"
      period      = "120"
      stat        = "p99"
      unit        = "Count"

      dimensions = {
        TargetGroup  = aws_alb_target_group.main.arn_suffix
        LoadBalancer = aws_alb.main.arn_suffix
      }
    }
  }
}

resource "aws_cloudwatch_metric_alarm" "httpcode_elb_5xx_count" {
  alarm_name          = "backend-load-balancer-5XX-${terraform.workspace}"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = "2"
  metric_name         = "HTTPCode_ELB_5XX_Count"
  namespace           = "AWS/ApplicationELB"
  period              = "60"
  statistic           = "Sum"
  threshold           = "25"
  alarm_description   = "This alarm indicates there are an abnormal amount of load balancer 5XX responses i.e it cannot talk with the Fleet backend target"
  alarm_actions       = [aws_sns_topic.cloudwatch_alarm_topic.arn]
  ok_actions          = [aws_sns_topic.cloudwatch_alarm_topic.arn]
  dimensions = {
    LoadBalancer = aws_alb.main.arn_suffix
  }
}

// Elasticache (redis) alerts https://docs.aws.amazon.com/AmazonElastiCache/latest/red-ug/CacheMetrics.WhichShouldIMonitor.html
resource "aws_cloudwatch_metric_alarm" "redis_cpu" {
  for_each            = toset(aws_elasticache_replication_group.default.member_clusters)
  alarm_name          = "redis-cpu-utilization-${each.key}-${terraform.workspace}"
  alarm_description   = "Redis cluster CPU utilization node ${each.key}"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = "1"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ElastiCache"
  period              = "300"
  statistic           = "Average"
  alarm_actions       = [aws_sns_topic.cloudwatch_alarm_topic.arn]
  ok_actions          = [aws_sns_topic.cloudwatch_alarm_topic.arn]

  threshold = "70"

  dimensions = {
    CacheClusterId = each.key
  }

  depends_on = [aws_elasticache_replication_group.default]
}

resource "aws_cloudwatch_metric_alarm" "redis_cpu_engine_utilization" {
  for_each            = toset(aws_elasticache_replication_group.default.member_clusters)
  alarm_name          = "redis-cpu-engine-utilization-${each.key}-${terraform.workspace}"
  alarm_description   = "Redis cluster CPU Engine utilization node ${each.key}"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = "1"
  metric_name         = "EngineCPUUtilization"
  namespace           = "AWS/ElastiCache"
  period              = "300"
  statistic           = "Average"
  alarm_actions       = [aws_sns_topic.cloudwatch_alarm_topic.arn]
  ok_actions          = [aws_sns_topic.cloudwatch_alarm_topic.arn]

  threshold = "25"

  dimensions = {
    CacheClusterId = each.key
  }

  depends_on = [aws_elasticache_replication_group.default]
}

resource "aws_cloudwatch_metric_alarm" "redis-current-connections" {
  for_each                  = toset(aws_elasticache_replication_group.default.member_clusters)
  alarm_name                = "redis-current-connections-${each.key}-${terraform.workspace}"
  alarm_description         = "Redis current connections for node ${each.key}"
  comparison_operator       = "LessThanLowerOrGreaterThanUpperThreshold"
  evaluation_periods        = "3"
  threshold_metric_id       = "e1"
  alarm_actions             = [aws_sns_topic.cloudwatch_alarm_topic.arn]
  ok_actions                = [aws_sns_topic.cloudwatch_alarm_topic.arn]
  insufficient_data_actions = []

  metric_query {
    id          = "e1"
    expression  = "ANOMALY_DETECTION_BAND(m1)"
    label       = "Current Connections (Expected)"
    return_data = "true"
  }

  metric_query {
    id          = "m1"
    return_data = "true"
    metric {
      metric_name = "CurrConnections"
      namespace   = "AWS/ElastiCache"
      period      = "300"
      stat        = "Average"
      unit        = "Count"

      dimensions = {
        CacheClusterId = each.key
      }
    }
  }
}

// ACM Certificate Manager
resource "aws_cloudwatch_metric_alarm" "acm_certificate_expired" {
  alarm_name          = "acm-cert-expiry-${terraform.workspace}"
  comparison_operator = "LessThanThreshold"
  evaluation_periods  = "1"
  period              = "86400" // 1 day in seconds
  statistic           = "Average"
  namespace           = "AWS/CertificateManager"
  metric_name         = "DaysToExpiry"
  actions_enabled     = "true"
  alarm_description   = "ACM Certificate will expire soon"
  alarm_actions       = [aws_sns_topic.cloudwatch_alarm_topic.arn]
  ok_actions          = [aws_sns_topic.cloudwatch_alarm_topic.arn]

  dimensions = {
    CertificateArn = aws_acm_certificate.dogfood_fleetdm_com.arn
  }
}