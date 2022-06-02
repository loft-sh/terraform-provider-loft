terraform {
  required_providers {
    loft = {
      source = "registry.terraform.io/loft-sh/loft"
    }
  }
}

provider "loft" {}

# Import data for all spaces
data "loft_spaces" "all" {
  cluster = "loft-cluster"
}

# Output all spaces
output "spaces" {
  value = data.loft_spaces.all.spaces.*.name
}
