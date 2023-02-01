terraform {
  required_providers {
    loft = {
      source = "registry.terraform.io/loft-sh/loft"
    }
  }
}

provider "loft" {}

resource "loft_space_instance" "example-space" {
  metadata {
    namespace = "loft-p-example-project"
    name      = "example-space"
  }
  spec {
    template_ref {
      name = "isolated-space"
    }
  }
}