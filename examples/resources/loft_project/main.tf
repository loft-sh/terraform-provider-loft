terraform {
  required_providers {
    loft = {
      source = "registry.terraform.io/loft-sh/loft"
    }
  }
}

provider "loft" {}

resource "loft_project" "example-project" {
  metadata {
    name = "example-project"
  }
  spec {
    access {
      name         = "loft-admin-access"
      verbs        = ["get", "update", "patch", "delete"]
      subresources = ["*"]
      users        = ["admin"]
    }
    access {
      name         = "loft-access"
      subresources = ["members", "clusters", "templates", "chartinfo", "charts"]
      verbs        = ["get"]
      users        = ["*"]
    }
    allowed_clusters {
      name = "*"
    }
    allowed_templates {
      kind  = "VirtualClusterTemplate"
      group = "storage.loft.sh"
      name  = "*"
    }
    allowed_templates {
      kind  = "SpaceTemplate"
      group = "storage.loft.sh"
      name  = "*"
    }
    description  = "Terraform Managed Project"
    display_name = "Terraform Managed Project"
    members {
      kind         = "User"
      group        = "storage.loft.sh"
      name         = "*"
      cluster_role = "loft-management-project-user"
    }
    owner {
      user = "admin"
    }
    quotas {
      project = {
        "spaceinstances"          = "10"
        "virtualclusterinstances" = "10"
      }
      user = {
        "spaceinstances"          = "10"
        "virtualclusterinstances" = "10"
      }
    }
  }
}