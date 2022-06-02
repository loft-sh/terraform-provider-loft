terraform {
  required_providers {
    loft = {
      source = "registry.terraform.io/loft-sh/loft"
    }
  }
}

provider "loft" {}

resource "loft_space" "sleep_after" {
  name            = "sleep-schedule-space"
  cluster         = "loft-cluster"
  sleep_schedule  = "* 18 * * *" # Sleep everyday at 6pm
  wakeup_schedule = "* 6 * * *"  # Wake everyday at 6am
}
