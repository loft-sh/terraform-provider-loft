terraform {
  required_providers {
    loft = {
      source  = "registry.terraform.io/loft-sh/loft"
      version = "0.0.1"
    }
  }
}

provider "loft" {}

resource "loft_space" "example" {
  name    = "example"
  cluster = "loft-cluster"
}

resource "loft_space" "dynamic" {
  generate_name = "example-"
  cluster       = "loft-cluster"
}

resource "loft_space" "sleepy_space" {
  name        = "sleepy"
  cluster     = "loft-cluster"
  sleep_after = 3600
  annotations = {
    "custom-annotation" = "special"
  }
}

data "loft_spaces" "all" {
  cluster = "loft-cluster"
}

output "spaces" {
  value = data.loft_spaces.all.spaces.*.name
}

output "example_name" {
  value = resource.loft_space.example.name
}

output "sleepy_space_annotations" {
  value = resource.loft_space.sleepy_space.annotations
}

output "sleepy_space_labels" {
  value = resource.loft_space.sleepy_space.labels
}