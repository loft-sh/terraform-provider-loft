terraform {
  required_providers {
    loft = {
      source = "registry.terraform.io/loft-sh/loft"
    }
  }
}

provider "loft" {}

resource "loft_space" "basic" {
  name    = "basic-space"
  cluster = "loft-cluster"
}

resource "loft_virtual_cluster" "basic" {
  name      = "basic-virtual-cluster"
  namespace = resource.loft_space.basic.name
  cluster   = resource.loft_space.basic.cluster
}
