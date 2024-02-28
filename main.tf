module "resource_group" {
  source   = "./resource_group"
  location = "eastus"
  name     = "avm-github-app"
}

module "acr" {
  source              = "./acr"
  location            = module.resource_group.resource_group_location
  registry_name       = "avmghapp"
  resource_group_name = module.resource_group.resource_group_name
  vnet_name           = "avmghapp"
}

module "docker_image" {
  source       = "./docker_image"
  password     = module.acr.push_password
  registry_url = module.acr.registry_url
  username     = module.acr.push_username
  image_tag    = var.image_tag
}

locals {
  port = 8080
}

module "container_apps" {
  source                         = "./container_apps"
  acr_url                        = module.acr.registry_url
  acr_user_name                  = module.acr.pull_username
  acr_user_password              = module.acr.pull_password
  container_app_environment_name = "avmghapp"
  docker_image                   = module.docker_image.docker_image
  location                       = module.resource_group.resource_group_location
  resource_group_name            = module.resource_group.resource_group_name
  subnet_id                      = module.acr.container_apps_subnet_id
  port                           = local.port
  github_app_config              = sensitive(yamlencode({
    server = {
      address = "0.0.0.0"
      port = local.port
    }
    expected_pusher_name = "azure-verified-module-draft[bot]"
    github = {
      v3_api_url = "https://api.github.com/"
      app = {
        integration_id = var.app_integration_id
        webhook_secret = var.webhook_secret
        private_key = base64decode(var.gh_app_private_key_pem_base64)
      }
    }
  }))
}