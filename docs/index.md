---
page_title: "Loft Provider"
subcategory: ""
description: |-
  Loft Provider (terraform-provider-loft)
---

# Loft Provider

The Loft Provider provides resources to manage your Loft Spaces and Virtual Clusters using Terraform.

~> This provider is deprecated for Loft versions 3.4 and above, and all versions of vCluster Platform. Users can now use the [kubernetes terraform provider](https://registry.terraform.io/providers/hashicorp/kubernetes/latest) with a [platform authenticated kube config](https://www.vcluster.com/docs/platform/api/authentication#log-in-via-cli).


## Example Usage

Create a Project using terraform
```terraform
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
```

Create a Space Instance using terraform
```terraform
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
```

Create a Virtual Cluster Instance using terraform
```terraform
terraform {
  required_providers {
    loft = {
      source = "registry.terraform.io/loft-sh/loft"
    }
  }
}

provider "loft" {}

resource "loft_virtual_cluster_instance" "example-vcluster" {
  metadata {
    namespace = "loft-p-example-project"
    name      = "example-vcluster"
  }
  spec {
    owner {
      user = "admin"
    }
    template_ref {
      name = "isolated-vcluster"
    }
  }
}
```

## Authentication and Configuration
Authentication for the Loft Provider can be provided through a Loft configuration file, or by manually providing the Loft host and access key.

### Loft Configuration File
When you login to Loft using the Loft CLI, a `config.json` file is create locally to store your login access key. By default, this Terraform provider will use this access key to authenticate when managing Loft resources. It may be neccessary to refresh your login using the [`loft login`](https://loft.sh/docs/cli/loft_login) command.

By default, the provider will authenticate using the currently logged in user:
```terraform
provider "loft" {
  # Uses the default Loft config location ($HOME/.loft/config.json)
}
```

To override the Loft config path location:
```terraform
provider "loft" {
  # If not using the default config location (`$HOME/.loft/config.json) you can change the location's `config_path`
  config_path = "/path/to/loft/config.json"
}
```

### Manual Configuration
The provider authentication can be manually configured using `access_key`, `host`, and `insecure` options. This is useful for when you want to configure authentication in a CI/CD environment and wish to provide credentials using secrets or environment variables.

This is an example using [terraform variables](https://www.terraform.io/language/values/variables) to set the `host`, `access_key`, and `insecure` options:
```terraform
variable "loft_host" {
}

variable "loft_access_key" {
}

variable "loft_insecure" {
}

provider "loft" {
  host       = var.loft_host
  access_key = var.loft_access_key
  insecure   = var.loft_insecure
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `access_key` (String) The Loft [access key](https://loft.sh/docs/api/access-keys).
- `config_path` (String) The Loft config file path. Defaults to `$HOME/.loft/config.json`.
- `host` (String) The Loft instance host.
- `insecure` (Boolean) Allow login into an insecure Loft instance. Defaults to `false`.