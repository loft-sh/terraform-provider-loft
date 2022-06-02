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
  sleep_after = 3600 # Sleep after 1 hour of inactivity
}

resource "loft_virtual_cluster" "vcluster_with_sleep_mode" {
  name      = "basic-vcluster"
  cluster   = resource.loft_space.sleep_after.cluster
  namespace = resource.loft_space.sleep_after.name
}
