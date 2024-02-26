variable "image_tag" {
  default = "latest"
}

variable "app_integration_id" {
  type = number
}

variable "webhook_secret" {
  type = string
  sensitive = true
}

variable "gh_app_private_key" {
  type = string
  sensitive = true
}