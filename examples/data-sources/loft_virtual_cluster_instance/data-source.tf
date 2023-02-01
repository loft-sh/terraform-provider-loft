terraform {
  required_providers {
    loft = {
      source = "registry.terraform.io/loft-sh/loft"
    }
  }
}

provider "loft" {}

# Output an existing virtual cluster instance not managed by terraform
data "loft_virtual_cluster_instance" "my-vcluster" {
  metadata {
    namespace = "loft-p-default"
    name      = "my-vcluster"
  }
}

output "vcluster" {
  value = data.loft_virtual_cluster_instance.my-vcluster.spec.0
}