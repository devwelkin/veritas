# 1. resource group
resource "azurerm_resource_group" "veritas_rg" {
  name     = "veritas-rg"
  location = "polandcentral"
}

# 2. container registry
resource "azurerm_container_registry" "veritas_acr" {
  name                = "veritasacr"
  resource_group_name = azurerm_resource_group.veritas_rg.name
  location            = azurerm_resource_group.veritas_rg.location
  sku                 = "Basic"
}

# 3. static public IP address for traefik
resource "azurerm_public_ip" "traefik_ip" {
  name                = "veritas-traefik-public-ip"
  resource_group_name = azurerm_resource_group.veritas_rg.name
  location            = azurerm_resource_group.veritas_rg.location
  allocation_method   = "Static"
  sku                 = "Standard"
}

# 4. kubernetes cluster
resource "azurerm_kubernetes_cluster" "veritas_aks" {
  name                = "veritas-aks"
  location            = azurerm_resource_group.veritas_rg.location
  resource_group_name = azurerm_resource_group.veritas_rg.name
  dns_prefix          = "veritas-aks" # since terraform creates it, we can give it a simple name
  sku_tier            = "Free"

  default_node_pool {
    name         = "default"
    node_count   = 1
    vm_size      = "Standard_B2s"
    os_disk_type = "Managed"
  }

  identity {
    type = "SystemAssigned"
  }
}

# 5. connection between aks and acr
resource "azurerm_role_assignment" "acr_pull_aks" {
  scope                = azurerm_container_registry.veritas_acr.id
  role_definition_name = "AcrPull"
  principal_id         = azurerm_kubernetes_cluster.veritas_aks.kubelet_identity[0].object_id
}



# 7. an output to easily get the ip address
output "static_public_ip" {
  value       = azurerm_public_ip.traefik_ip.ip_address
  description = "The static public IP address for the Traefik Ingress."
}