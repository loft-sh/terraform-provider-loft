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