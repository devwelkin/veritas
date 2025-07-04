terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.11"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.23"
    }
  }
}

provider "azurerm" {
  features {}
}

provider "helm" {
  kubernetes {
    host                   = azurerm_kubernetes_cluster.veritas_aks.kube_config.0.host
    client_certificate     = base64decode(azurerm_kubernetes_cluster.veritas_aks.kube_config.0.client_certificate)
    client_key             = base64decode(azurerm_kubernetes_cluster.veritas_aks.kube_config.0.client_key)
    cluster_ca_certificate = base64decode(azurerm_kubernetes_cluster.veritas_aks.kube_config.0.cluster_ca_certificate)
  }
}

provider "kubernetes" {
  host                   = azurerm_kubernetes_cluster.veritas_aks.kube_config.0.host
  client_certificate     = base64decode(azurerm_kubernetes_cluster.veritas_aks.kube_config.0.client_certificate)
  client_key             = base64decode(azurerm_kubernetes_cluster.veritas_aks.kube_config.0.client_key)
  cluster_ca_certificate = base64decode(azurerm_kubernetes_cluster.veritas_aks.kube_config.0.cluster_ca_certificate)
}