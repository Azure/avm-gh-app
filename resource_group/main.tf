resource "azurerm_resource_group" "this" {
  location = var.location
  name     = var.name

  tags = {
    env = "prod"
  }
}