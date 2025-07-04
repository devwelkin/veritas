# specifies the cloud provider and version for terraform
terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.0" # approximately version 3.0
    }
  }
}

# settings for connecting to azure
# this block will remain empty because we have already authenticated with 'az login'
provider "azurerm" {
  features {}
}

# --- resources ---

# 1. resource group (az group create...)
resource "azurerm_resource_group" "veritas_rg" {
  name     = "veritas-rg"
  location = "polandcentral" 
}

# 2. container registry (az acr create...)
resource "azurerm_container_registry" "veritas_acr" {
  name                = "veritasacr" # should be the same name you gave in azure
  resource_group_name = azurerm_resource_group.veritas_rg.name
  location            = azurerm_resource_group.veritas_rg.location
  sku                 = "Basic"
  admin_enabled       = false # better to be disabled for security
}

# 3. kubernetes cluster (az aks create...)
resource "azurerm_kubernetes_cluster" "veritas_aks" {
  name                = "veritas-aks"
  # your plan shows the node resource group is in polandcentral, not westeurope
  location            = "polandcentral" 
  resource_group_name = azurerm_resource_group.veritas_rg.name
  # use the *actual* dns prefix that azure created
  dns_prefix          = "veritas-ak-veritas-rg-ff15cd"
  sku_tier            = "Free"

  default_node_pool {
    # use the *actual* node pool name
    name       = "nodepool1" 
    node_count = 1
    # use the *actual* vm size
    vm_size    = "Standard_B2s"
    # match the os disk type
    os_disk_type = "Managed"
  }

  identity {
    type = "SystemAssigned"
  }
  
  # --- ADD THIS BLOCK ---
  # this tells terraform to not worry if these specific parts
  # of the live resource don't match the code.
  lifecycle {
    ignore_changes = [
      linux_profile,
      # it's also a good idea to ignore tags, as azure sometimes adds its own
      tags, 
    ]
  }
}

# 4. connect aks and acr (az aks update --attach-acr...)
# this creates the necessary role assignment for aks to be able to pull images from acr
resource "azurerm_role_assignment" "acr_pull_aks" {
  scope                = azurerm_container_registry.veritas_acr.id
  role_definition_name = "AcrPull"
  principal_id         = azurerm_kubernetes_cluster.veritas_aks.kubelet_identity[0].object_id
}