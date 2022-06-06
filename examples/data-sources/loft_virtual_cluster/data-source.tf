terraform {
  required_providers {
    loft = {
      source = "registry.terraform.io/loft-sh/loft"
    }
  }
}

provider "loft" {}

# Import an existing virtual cluster not managed by terraform
data "loft_virtual_cluster" "existing_virtual_cluster" {
  cluster   = "loft-cluster"
  namespace = "default"
  name      = "existing-virtual-cluster"
}
