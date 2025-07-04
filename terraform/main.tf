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

# 6. install traefik helm chart
resource "helm_release" "traefik" {
  name       = "traefik"
  repository = "https://helm.traefik.io/traefik"
  chart      = "traefik"
  version    = "25.0.0"

  
  timeout = 600

  depends_on = [azurerm_role_assignment.acr_pull_aks]
}
