locals {
  prefix = "bioinformatics-${var.env}-spark"

  common_tags = {
    "env"       = var.env
    "role"      = "spark"
    "terraform" = "true"
  }
}
