terraform {
  required_version = ">=1.0,<2.0"

  required_providers {
    kubernetes = ">=2.21.1,<3"
  }
}

provider "kubernetes" {
  config_path    = "~/.kube/config"
  config_context = "loft_loft-cluster"
  insecure = true
}
#
#resource "kubernetes_manifest" "my_secret" {
#  manifest = {
#    "apiVersion" = "management.loft.sh/v1"
#    "kind"       = "SharedSecret"
#    "metadata" = {
#      "name"      = "my-secret"
#      "namespace" = "loft"
#    }
#    "spec" = {
#      "description" = "Shared Secret Demo"
#      "displayName" = "Shared Secret Demo"
#      "data" = {
#        "name"         = base64encode("Demo")
#      }
#      "access" = [
#        {
#          "verbs"        = ["*"]
#          "subresources" = ["*"]
#          "teams"        = ["loft-admins"]
#        }
#      ]
#    }
#  }
#}