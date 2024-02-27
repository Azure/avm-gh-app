variable "location" {
  type = string
}

variable "resource_group_name" {
  type = string
}

variable "subnet_id" {
  type = string
}

variable "acr_url" {
  type = string
}

variable "acr_user_name" {
  type = string
}

variable "acr_user_password" {
  type      = string
  sensitive = true
}

variable "docker_image" {
  type = string
}

variable "container_app_environment_name" {
  type = string
}

variable "port" {
  type = number
}