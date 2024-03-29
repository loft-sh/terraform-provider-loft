---
page_title: "loft_project Resource - terraform-provider-loft"
subcategory: ""
description: |-
Provides details for loft_project Resource
---
# loft_project (Resource)
Project holds the Project information

## Example Usage
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

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `metadata` (Block List, Min: 1, Max: 1) Standard Project's metadata. More info: https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#metadata (see [below for nested schema](#nestedblock--metadata))
- `spec` (Block List, Min: 1) (see [below for nested schema](#nestedblock--spec))

### Read-Only

- `id` (String) Unique identifier for this resource. The format is `<name>`.

<a id="nestedblock--metadata"></a>
### Nested Schema for `metadata`

Optional:

- `annotations` (Map of String) An unstructured key value map stored with the Project that may be used to store arbitrary metadata. More info: http://kubernetes.io/docs/user-guide/annotations
- `generate_name` (String) Prefix, used by the server, to generate a unique name ONLY IF the `name` field has not been provided. This value will also be combined with a unique suffix. Read more: https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#idempotency
- `labels` (Map of String) Map of string keys and values that can be used to organize and categorize (scope and select) the Project. May match selectors of replication controllers and services. More info: http://kubernetes.io/docs/user-guide/labels
- `name` (String) Name of the Project, must be unique. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/identifiers#names

Read-Only:

- `generation` (Number) A sequence number representing a specific generation of the desired state.
- `resource_version` (String) An opaque value that represents the internal version of this Project that can be used by clients to determine when Project has changed. Read more: https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#concurrency-control-and-consistency
- `uid` (String) The unique in time and space value for this Project. More info: http://kubernetes.io/docs/user-guide/identifiers#uids


<a id="nestedblock--spec"></a>
### Nested Schema for `spec`

Optional:

- `access` (Block List) Access holds the access rights for users and teams (see [below for nested schema](#nestedblock--spec--access))
- `allowed_clusters` (Block List) AllowedClusters are target clusters that are allowed to target with environments. (see [below for nested schema](#nestedblock--spec--allowed_clusters))
- `allowed_templates` (Block List) AllowedTemplates are the templates that are allowed to use in this project. (see [below for nested schema](#nestedblock--spec--allowed_templates))
- `argo_c_d` (Block List, Max: 1) ArgoIntegration holds information about ArgoCD Integration (see [below for nested schema](#nestedblock--spec--argo_c_d))
- `description` (String) Description describes an app
- `display_name` (String) DisplayName is the name that should be displayed in the UI
- `members` (Block List) Members are the users and teams that are part of this project (see [below for nested schema](#nestedblock--spec--members))
- `namespace_pattern` (Block List, Max: 1) NamespacePattern specifies template patterns to use for creating each space or virtual cluster's namespace (see [below for nested schema](#nestedblock--spec--namespace_pattern))
- `owner` (Block List, Max: 1) Owner holds the owner of this object (see [below for nested schema](#nestedblock--spec--owner))
- `quotas` (Block List, Max: 1) Quotas define the quotas inside the project (see [below for nested schema](#nestedblock--spec--quotas))

<a id="nestedblock--spec--access"></a>
### Nested Schema for `spec.access`

Required:

- `verbs` (List of String) Verbs is a list of Verbs that apply to ALL the ResourceKinds and AttributeRestrictions contained in this rule. VerbAll represents all kinds.

Optional:

- `name` (String) Name is an optional name that is used for this access rule
- `subresources` (List of String) Subresources defines the sub resources that are allowed by this access rule
- `teams` (List of String) Teams specifies which teams should be able to access this secret with the aforementioned verbs
- `users` (List of String) Users specifies which users should be able to access this secret with the aforementioned verbs


<a id="nestedblock--spec--allowed_clusters"></a>
### Nested Schema for `spec.allowed_clusters`

Optional:

- `name` (String) Name is the name of the cluster that is allowed to create an environment in


<a id="nestedblock--spec--allowed_templates"></a>
### Nested Schema for `spec.allowed_templates`

Optional:

- `group` (String) Group of the template that is allowed. Currently only supports storage.loft.sh
- `is_default` (Boolean) IsDefault specifies if the template should be used as a default
- `kind` (String) Kind of the template that is allowed. Currently only supports VirtualClusterTemplate & SpaceTemplate
- `name` (String) Name of the template


<a id="nestedblock--spec--argo_c_d"></a>
### Nested Schema for `spec.argo_c_d`

Optional:

- `cluster` (String) Cluster defines the name of the cluster that ArgoCD is deployed into -- if not provided this will default to 'loft-cluster'.
- `enabled` (Boolean) Enabled indicates if the ArgoCD Integration is enabled for the project -- this knob only enables the syncing of virtualclusters, but does not enable SSO integration or project creation (see subsequent spec sections!).
- `namespace` (String) Namespace defines the namespace in which ArgoCD is running in the cluster.
- `project` (Block List, Max: 1) Project defines project related values for the ArgoCD Integration. Enabling Project integration will cause Loft to generate and manage an ArgoCD appProject that corresponds to the Loft Project. (see [below for nested schema](#nestedblock--spec--argo_c_d--project))
- `sso` (Block List, Max: 1) SSO defines single-sign-on related values for the ArgoCD Integration. Enabling SSO will allow users to authenticate to ArgoCD via Loft. (see [below for nested schema](#nestedblock--spec--argo_c_d--sso))
- `virtual_cluster_instance` (String) VirtualClusterInstance defines the name of *virtual cluster* (instance) that ArgoCD is deployed into. If provided, Cluster will be ignored and Loft will assume that ArgoCD is running in the specified virtual cluster.

<a id="nestedblock--spec--argo_c_d--project"></a>
### Nested Schema for `spec.argo_c_d.project`

Optional:

- `enabled` (Boolean) Enabled indicates if the ArgoCD Project Integration is enabled for this project. Enabling this will cause Loft to create an appProject in ArgoCD that is associated with the Loft Project. When Project integration is enabled Loft will override the default assigned role set in the SSO integration spec.
- `metadata` (Block List, Max: 1) Metadata defines additional metadata to attach to the loft created project in ArgoCD. (see [below for nested schema](#nestedblock--spec--argo_c_d--project--metadata))
- `roles` (Block List) Roles is a list of roles that should be attached to the ArgoCD project. If roles are provided no loft default roles will be set. If no roles are provided *and* SSO is enabled, loft will configure sane default values. (see [below for nested schema](#nestedblock--spec--argo_c_d--project--roles))
- `source_repos` (List of String) SourceRepos is a list of source repositories to attach/allow on the project, if not specified will be "*" indicating all source repositories.

<a id="nestedblock--spec--argo_c_d--project--metadata"></a>
### Nested Schema for `spec.argo_c_d.project.metadata`

Optional:

- `description` (String) Description to add to the ArgoCD project.
- `extra_annotations` (Map of String) ExtraAnnotations are optional annotations that can be attached to the project in ArgoCD.
- `extra_labels` (Map of String) ExtraLabels are optional labels that can be attached to the project in ArgoCD.


<a id="nestedblock--spec--argo_c_d--project--roles"></a>
### Nested Schema for `spec.argo_c_d.project.roles`

Optional:

- `description` (String) Description to add to the ArgoCD project.
- `groups` (List of String) Groups is a list of OIDC group names to bind to the role.
- `name` (String) Name of the ArgoCD role to attach to the project.
- `rules` (Block List) Rules ist a list of policy rules to attach to the role. (see [below for nested schema](#nestedblock--spec--argo_c_d--project--roles--rules))

<a id="nestedblock--spec--argo_c_d--project--roles--rules"></a>
### Nested Schema for `spec.argo_c_d.project.roles.rules`

Optional:

- `action` (String) Action is one of "*", "get", "create", "update", "delete", "sync", or "override".
- `application` (String) Application is the ArgoCD project/repository to apply the rule to.
- `permission` (Boolean) Allow applies the "allow" permission to the rule, if allow is not set, the permission will always be set to "deny".




<a id="nestedblock--spec--argo_c_d--sso"></a>
### Nested Schema for `spec.argo_c_d.sso`

Optional:

- `assigned_roles` (List of String) AssignedRoles is a list of roles to assign for users who authenticate via Loft -- by default this will be the `read-only` role. If any roles are provided this will override the default setting.
- `enabled` (Boolean) Enabled indicates if the ArgoCD SSO Integration is enabled for this project. Enabling this will cause Loft to configure SSO authentication via Loft in ArgoCD. If Projects are *not* enabled, all users associated with this Project will be assigned either the 'read-only' (default) role, *or* the roles set under the AssignedRoles field.
- `host` (String) Host defines the ArgoCD host address that will be used for OIDC authentication between loft and ArgoCD. If not specified OIDC integration will be skipped, but vclusters/spaces will still be synced to ArgoCD.



<a id="nestedblock--spec--members"></a>
### Nested Schema for `spec.members`

Optional:

- `cluster_role` (String) ClusterRole is the assigned role for the above member
- `group` (String) Group of the member. Currently only supports storage.loft.sh
- `kind` (String) Kind is the kind of the member. Currently either User or Team
- `name` (String) Name of the member


<a id="nestedblock--spec--namespace_pattern"></a>
### Nested Schema for `spec.namespace_pattern`

Optional:

- `space` (String) Space holds the namespace pattern to use for space instances
- `virtual_cluster` (String) VirtualCluster holds the namespace pattern to use for virtual cluster instances


<a id="nestedblock--spec--owner"></a>
### Nested Schema for `spec.owner`

Optional:

- `team` (String) Team specifies a Loft team.
- `user` (String) User specifies a Loft user.


<a id="nestedblock--spec--quotas"></a>
### Nested Schema for `spec.quotas`

Optional:

- `project` (Map of String) Project holds the quotas for the whole project
- `user` (Map of String) User holds the quotas per user / team

## Import
Import is supported using the following syntax:
```shell
# import the `example-project` into the `loft_project.example-project` resource
terraform import loft_project.example-project example-project
```