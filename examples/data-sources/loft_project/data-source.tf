terraform {
  required_providers {
    loft = {
      source = "registry.terraform.io/loft-sh/loft"
    }
  }
}

provider "loft" {}

# Output an existing project not managed by terraform
data "loft_project" "default" {
  metadata {
    name = "default"
  }
}

output "project" {
  value = data.loft_project.default.spec.0
}