variable "image_tag" {
  default = "latest"
}

variable "gh_app_integration_id" {
  type = number
}

variable "gh_app_webhook_secret" {
  type = string
  sensitive = true
}

# cat key.pem | base64 -w 0
variable "gh_app_private_key_pem_base64" {
  type = string
  sensitive = true
}