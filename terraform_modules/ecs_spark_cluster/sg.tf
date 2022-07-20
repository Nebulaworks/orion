resource "aws_security_group" "spark_cluster" {
  name        = "${local.prefix}-cluster"
  vpc_id      = var.vpc_id
  description = "Security group for the spark cluster"

  ingress {
    description     = "Allows connectivity to the spark master for job execution"
    from_port       = 7077
    to_port         = 7077
    protocol        = "tcp"
    security_groups = [aws_security_group.spark_nlb.id]
  }

  ingress {
    description     = "Allows connectivity to the spark master to access the UI"
    from_port       = 8080
    to_port         = 8080
    protocol        = "tcp"
    security_groups = [aws_security_group.spark_alb.id]
  }

  egress {
    description = "General Egress"
    protocol    = "-1"
    from_port   = 0
    to_port     = 0
    cidr_blocks = ["0.0.0.0/0"] #tfsec:ignore:aws-vpc-no-public-egress-sgr
  }

  tags = merge({ "Name" = "${local.prefix}-cluster" }, local.common_tags)
}

resource "aws_security_group" "spark_alb" {
  name        = "${local.prefix}-alb"
  vpc_id      = var.vpc_id
  description = "Security group for the spark ALB"

  ingress {
    description     = "Allows connectivity to the spark master to access the UI"
    from_port       = 8080
    to_port         = 8080
    protocol        = "tcp"
    security_groups = var.security_group_whitelist
  }

  egress {
    description     = "Egress to Spark Cluster"
    protocol        = "tcp"
    from_port       = 8080
    to_port         = 8080
    security_groups = [aws_security_group.spark_cluster.id]
  }

  tags = merge({ "Name" = "${local.prefix}-alb" }, local.common_tags)
}

resource "aws_security_group" "spark_nlb" {
  name        = "${local.prefix}-nlb"
  vpc_id      = var.vpc_id
  description = "Security group for the spark NLB"

  ingress {
    description     = "Allows connectivity to the spark master to access the UI"
    from_port       = 7077
    to_port         = 7077
    protocol        = "tcp"
    security_groups = var.security_group_whitelist
  }

  egress {
    description     = "Egress to Spark Cluster"
    protocol        = "tcp"
    from_port       = 7077
    to_port         = 7077
    security_groups = [aws_security_group.spark_cluster.id]
  }

  tags = merge({ "Name" = "${local.prefix}-nlb" }, local.common_tags)
}
