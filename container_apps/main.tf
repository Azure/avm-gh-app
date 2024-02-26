locals {
  port = var.port
}

resource "azurerm_application_insights" "this" {
  application_type    = "other"
  location            = var.location
  name                = "avmghapp"
  resource_group_name = var.resource_group_name
}

module "avmgithubapp" {
  source                                             = "Azure/container-apps/azure"
  version                                            = "0.2.0"
  container_app_environment_name                     = var.container_app_environment_name
  container_app_environment_infrastructure_subnet_id = var.subnet_id
  container_apps = {
    avmghapp = {
      name          = "avmghapp"
      revision_mode = "Single"
      registry = [
        {
          server               = var.acr_url
          username             = var.acr_user_name
          password_secret_name = "passsec"
        }
      ]
      template = {
        containers = [
          {
            name   = "avmghapp"
            memory = "0.5Gi"
            cpu    = 0.25
            image  = var.docker_image
            env = toset([
              {
                name  = "GITHUB_APP_CONFIG"
                secret_name = "config"
              },
              ])
          }
        ]
      }
      ingress = {
        allow_insecure_connection = false
        external_enabled          = true
        target_port               = local.port
        traffic_weight = {
          latest_revision = true
          percentage      = 100
        }
      }
    }
  }
  container_app_secrets = {
    avmghapp = [
      {
        name  = "passsec"
        value = var.acr_user_password
      },
      {
        name  = "config"
        value = var.github_app_config
      }
    ]
  }
  location                     = var.location
  log_analytics_workspace_name = "avmghapp-log-analytics-workspace"
  resource_group_name          = var.resource_group_name
}