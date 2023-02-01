terraform {
  required_providers {
    loft = {
      source = "registry.terraform.io/loft-sh/loft"
    }
  }
}

provider "loft" {}

# Output an existing space instance not managed by terraform
data "loft_space_instance" "my-space" {
  metadata {
    namespace = "loft-p-default"
    name      = "my-space"
  }
}

output "space" {
  value = data.loft_space_instance.my-space.spec.0
}