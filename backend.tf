terraform {
  backend "azurerm" {
    snapshot             = true
    use_msi              = true
  }
}