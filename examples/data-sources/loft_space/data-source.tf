terraform {
  required_providers {
    loft = {
      source = "registry.terraform.io/loft-sh/loft"
    }
  }
}

provider "loft" {}

# Import an existing space not managed by terraform
data "loft_space" "default" {
  cluster = "loft-cluster"
  name    = "default"
}
