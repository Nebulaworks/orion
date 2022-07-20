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

variable "cluster_subnets" {
  description = "The subnets in which to deploy the services. Required."
  type        = list(string)
}  

variable "spark_version" {
  description = "The version of Apache Spark to run. Determines which container to pull down. Required."
  type        = string
}

variable "container_sha" {
  description = "The SHA of the container to which to pin the task definition."
  type        = string
  default     = "26f697fd5df82abf1eb71ccff6ac516ec5fcf476c3ffb9bed3f1668a210dcb8f"
}

variable "master_cpu" {
  description = "The amount of CPU to dedicate to the spark_master service (needs to be a multiple of 1024)."
  type        = number
  default     = 2048
}

variable "master_memory" {
  description = "The amount of memory to dedicate to the spark_master service (needs to be a multiple of 1024)."
  type        = number
  default     = 4096
}

variable "worker_cpu" {
  description = "The amount of CPU to dedicate to the spark_master service (needs to be a multiple of 1024)."
  type        = number
  default     = 2048
}

variable "worker_memory" {
  description = "The amount of memory to dedicate to the spark_master service (needs to be a multiple of 1024)."
  type        = number
  default     = 4096
}

variable "worker_count" {
  description = "The number of spark workers to keep in the cluster"
  type        = number
  default     = 3
}
