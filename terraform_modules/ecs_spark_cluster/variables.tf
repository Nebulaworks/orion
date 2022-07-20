variable "env" {
  description = "The environment prefix that this resource is associated with."
  type        = string
  default     = "dev"
}

variable "role" {
  description = "The application role that will be utilizing the spark cluster resources. Required."
  type        = string
}

variable "region" {
  description = "The region in which the cluster and associated resources will reside. Required."
  type        = string
}

variable "vpc_id" {
  description = "The ID of the VPC containing the spark cluster and associated resources. Required."
  type        = string
}
