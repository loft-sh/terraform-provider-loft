terraform {
  required_providers {
    loft = {
      source = "registry.terraform.io/loft-sh/loft"
    }
  }
}

provider "loft" {}

resource "loft_virtual_cluster_instance" "example-vcluster" {
  metadata {
    namespace = "loft-p-example-project"
    name      = "example-vcluster"
  }
  spec {
    owner {
      user = "admin"
    }
    template_ref {
      name = "isolated-vcluster"
    }
  }
}