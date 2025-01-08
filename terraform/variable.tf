variable "region" {
  description = "The AWS region where the resources will be created"
  type        = string
  default     = "eu-west-1"
}

variable "docker_image_name" {
  description = "The name of the Docker image to build and push to ECR"
  type        = string
  default     = "go-ticket-api"
}

variable "docker_image_tag" {
  description = "The tag of the Docker image to build and push to ECR"
  type        = string
  default     = "latest"
}
