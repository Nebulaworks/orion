locals {
  prefix = "${var.role}-${var.env}-spark"

  common_tags = {
    "env"       = var.env
    "role"      = var.role
    "terraform" = "true"
  }
}
