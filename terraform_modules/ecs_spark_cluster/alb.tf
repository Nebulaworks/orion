# ALB for Spark UI
resource "aws_lb" "spark_ui" {
  name                       = "${local.prefix}-ui"
  internal                   = true
  load_balancer_type         = "application"
  drop_invalid_header_fields = true

  security_groups = [aws_security_group.spark_alb.id]
  subnets         = var.cluster_subnets

  tags = merge({ "Name" = "${local.prefix}-ui" }, local.common_tags)
}

resource "aws_lb_listener" "spark_ui_http" {
  load_balancer_arn = aws_lb.spark_ui.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type = "redirect"

    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
    }
  }
}

resource "aws_lb_listener" "spark_ui_https" {
  load_balancer_arn = aws_lb.spark_ui.arn
  port              = 443
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-TLS-1-2-2017-01"
  certificate_arn   = aws_acm_certificate.spark_ui.arn

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.spark_ui.arn
  }
}

resource "aws_lb_target_group" "spark_ui" {
  name = "${local.prefix}-ui"

  vpc_id      = var.vpc_id
  port        = 8080
  protocol    = "HTTP"
  target_type = "ip"

  health_check {
    enabled = true
    path    = "/"
    matcher = "200"
  }

  tags = merge({ "Name" = "${local.prefix}-ui" }, local.common_tags)
}

# NLB for Spark Jobs
resource "aws_lb" "spark_jobs" {
  name               = "${local.prefix}-jobs"
  internal           = true
  load_balancer_type = "network"

  security_groups = [aws_security_group.spark_nlb.id]
  subnets         = var.cluster_subnets

  tags = merge({ "Name" = "${local.prefix}-jobs" }, local.common_tags)
}

resource "aws_lb_listener" "spark_jobs" {
  load_balancer_arn = aws_lb.spark_jobs.arn
  port              = 7077
  protocol          = "TCP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.spark_jobs.arn
  }
}

resource "aws_lb_target_group" "spark_jobs" {
  name = "${local.prefix}-jobs"

  vpc_id      = var.vpc_id
  port        = 7077
  protocol    = "TCP"
  target_type = "ip"

  health_check {
    enabled = true
    port    = "traffic-port"
  }

  tags = merge({ "Name" = "${local.prefix}-jobs" }, local.common_tags)
}

