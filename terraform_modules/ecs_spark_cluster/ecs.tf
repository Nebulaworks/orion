resource "aws_ecs_cluster" "spark" {
  name = local.prefix

  tags = merge({ "Name" = "${local.prefix}" }, local.common_tags)
}

resource "aws_ecs_task_definition" "spark_master" {
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = var.master_cpu
  memory                   = var.master_memory
  family                   = local.prefix
  execution_role_arn       = aws_iam_role.task_execution_role.arn

  container_definitions = jsonencode([{
    name  = "${local.prefix}-master"
    image = "docker.io/bitnami/spark@sha256:${var.container_sha}"
    environment = [{
      name  = "SPARK_MODE"
      value = "master"
      },
      {
        name  = "SPARK_RPC_AUTHENTICATION_ENABLED"
        value = "no"
      },
      {
        name  = "SPARK_RPC_ENCRYPTION_ENABLED"
        value = "no"
      },
      {
        name  = "SPARK_LOCAL_STORAGE_ENCRYPTION_ENABLED"
        value = "no"
      },
      {
        name  = "SPARK_SSL_ENABLED"
        value = "no"
    }]
    essential = true
    portMappings = [{
      protocol      = "tcp"
      containerPort = 8080
      hostPort      = 8080
      },
      {
        protocol      = "tcp"
        containerPort = 7077
        hostPort      = 7077
    }]
  }])

  tags = merge({ "Name" = "${local.prefix}" }, local.common_tags)
}

resource "aws_ecs_service" "spark_master" {
  name                               = local.prefix
  cluster                            = aws_ecs_cluster.spark.id
  task_definition                    = aws_ecs_task_definition.spark_master.arn
  desired_count                      = 1
  deployment_minimum_healthy_percent = 0
  deployment_maximum_percent         = 100
  launch_type                        = "FARGATE"
  scheduling_strategy                = "REPLICA"

  load_balancer {
    container_name   = "termapply"
    container_port   = "7077"
    target_group_arn = aws_lb_target_group.spark_job.arn
  }

  load_balancer {
    container_name   = "termapply"
    container_port   = "8080"
    target_group_arn = aws_lb_target_group.spark_ui.arn
  }

  network_configuration {
    security_groups = [aws_security_group.spark.id]
    subnets         = var.cluster_subnets
  }

  tags = merge({ "Name" = "${local.prefix}" }, local.common_tags)
}

resource "aws_ecs_task_definition" "spark_worker" {
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = var.worker_cpu
  memory                   = var.worker_memory
  family                   = local.prefix
  execution_role_arn       = aws_iam_role.task_execution_role.arn

  container_definitions = jsonencode([{
    name  = "${local.prefix}-worker"
    image = "docker.io/bitnami/spark@sha256:${var.container_sha}"
    environment = [{
      name  = "SPARK_MODE"
      value = "worker"
      },
      {
        name  = "SPARK_MASTER_URL"
        value = "spark://spark:7077"
      },
      {
        name  = "SPARK_WORKER_MEMORY"
        value = "${var.worker_memory / 1024}G"
      },
      {
        name  = "SPARK_WORKER_CORES"
        value = "${var.worker_cpu / 1024}"
      },
      {
        name  = "SPARK_RPC_AUTHENTICATION_ENABLED"
        value = "no"
      },
      {
        name  = "SPARK_RPC_ENCRYPTION_ENABLED"
        value = "no"
      },
      {
        name  = "SPARK_LOCAL_STORAGE_ENCRYPTION_ENABLED"
        value = "no"
      },
      {
        name  = "SPARK_SSL_ENABLED"
        value = "no"
    }]
    essential = true
    portMappings = [{
      protocol      = "tcp"
      containerPort = 8080
      hostPort      = 8080
      },
      {
        protocol      = "tcp"
        containerPort = 7077
        hostPort      = 7077
    }]
  }])

  tags = merge({ "Name" = "${local.prefix}" }, local.common_tags)
}

resource "aws_ecs_service" "spark_worker" {
  name                               = local.prefix
  cluster                            = aws_ecs_cluster.spark.id
  task_definition                    = aws_ecs_task_definition.spark_worker.arn
  desired_count                      = var.worker_count
  deployment_minimum_healthy_percent = 0
  deployment_maximum_percent         = 100
  launch_type                        = "FARGATE"
  scheduling_strategy                = "REPLICA"

  network_configuration {
    security_groups = [aws_security_group.spark.id]
    subnets         = var.cluster_subnets
  }

  tags = merge({ "Name" = "${local.prefix}" }, local.common_tags)
}
