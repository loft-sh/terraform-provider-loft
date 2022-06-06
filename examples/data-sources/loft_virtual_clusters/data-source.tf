terraform {
  required_providers {
    loft = {
      source = "registry.terraform.io/loft-sh/loft"
    }
  }
}

provider "loft" {}

# Import data for all virtual clusters in the 'default' namespace
data "loft_virtual_clusters" "all" {
  cluster   = "loft-cluster"
  namespace = "default"
}

# Output all virtual clusters
output "all_virtual_clusters" {
  value = data.loft_virtual_clusters.all.virtual_clusters.*.name
}
