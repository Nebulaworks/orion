resource "aws_iam_role" "task_execution_role" {
  name_prefix        = "${local.prefix}-task"
  assume_role_policy = data.aws_iam_policy_document.ecs_assumerole.json
}

data "aws_iam_policy_document" "ecs_assumerole" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }
  }
}

resource "aws_iam_role_policy_attachment" "ecs_task_execution" {
  role       = aws_iam_role.task_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}
