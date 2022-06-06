terraform {
  required_providers {
    loft = {
      source = "registry.terraform.io/loft-sh/loft"
    }
  }
}

provider "loft" {}

resource "loft_space" "sleep_after" {
  name        = "sleep-mode-space"
  cluster     = "loft-cluster"
  sleep_after = "18000s"
}