---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "loft_project Data Source - terraform-provider-loft"
subcategory: ""
description: |-
  Provides details for loft_project Data Source
---

# loft_project (Data Source)

Project holds the Project information



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `metadata` (Block List, Min: 1, Max: 1) Standard Project's metadata. More info: https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#metadata (see [below for nested schema](#nestedblock--metadata))

### Read-Only

- `id` (String) Unique identifier for this resource. The format is `<name>`.
- `spec` (List of Object) (see [below for nested schema](#nestedatt--spec))

<a id="nestedblock--metadata"></a>
### Nested Schema for `metadata`

Required:

- `name` (String) Name of the Project, must be unique. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/identifiers#names

Optional:

- `annotations` (Map of String) An unstructured key value map stored with the Project that may be used to store arbitrary metadata. More info: http://kubernetes.io/docs/user-guide/annotations
- `generate_name` (String) Prefix, used by the server, to generate a unique name ONLY IF the `name` field has not been provided. This value will also be combined with a unique suffix. Read more: https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#idempotency
- `labels` (Map of String) Map of string keys and values that can be used to organize and categorize (scope and select) the Project. May match selectors of replication controllers and services. More info: http://kubernetes.io/docs/user-guide/labels

Read-Only:

- `generation` (Number) A sequence number representing a specific generation of the desired state.
- `resource_version` (String) An opaque value that represents the internal version of this Project that can be used by clients to determine when Project has changed. Read more: https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#concurrency-control-and-consistency
- `uid` (String) The unique in time and space value for this Project. More info: http://kubernetes.io/docs/user-guide/identifiers#uids


<a id="nestedatt--spec"></a>
### Nested Schema for `spec`

Read-Only:

- `access` (List of Object) (see [below for nested schema](#nestedobjatt--spec--access))
- `allowed_clusters` (List of Object) (see [below for nested schema](#nestedobjatt--spec--allowed_clusters))
- `allowed_templates` (List of Object) (see [below for nested schema](#nestedobjatt--spec--allowed_templates))
- `argo_c_d` (List of Object) (see [below for nested schema](#nestedobjatt--spec--argo_c_d))
- `description` (String)
- `display_name` (String)
- `members` (List of Object) (see [below for nested schema](#nestedobjatt--spec--members))
- `namespace_pattern` (List of Object) (see [below for nested schema](#nestedobjatt--spec--namespace_pattern))
- `owner` (List of Object) (see [below for nested schema](#nestedobjatt--spec--owner))
- `quotas` (List of Object) (see [below for nested schema](#nestedobjatt--spec--quotas))

<a id="nestedobjatt--spec--access"></a>
### Nested Schema for `spec.access`

Read-Only:

- `name` (String)
- `subresources` (List of String)
- `teams` (List of String)
- `users` (List of String)
- `verbs` (List of String)


<a id="nestedobjatt--spec--allowed_clusters"></a>
### Nested Schema for `spec.allowed_clusters`

Read-Only:

- `name` (String)


<a id="nestedobjatt--spec--allowed_templates"></a>
### Nested Schema for `spec.allowed_templates`

Read-Only:

- `group` (String)
- `is_default` (Boolean)
- `kind` (String)
- `name` (String)


<a id="nestedobjatt--spec--argo_c_d"></a>
### Nested Schema for `spec.argo_c_d`

Read-Only:

- `cluster` (String)
- `enabled` (Boolean)
- `namespace` (String)
- `project` (List of Object) (see [below for nested schema](#nestedobjatt--spec--argo_c_d--project))
- `sso` (List of Object) (see [below for nested schema](#nestedobjatt--spec--argo_c_d--sso))
- `virtual_cluster_instance` (String)

<a id="nestedobjatt--spec--argo_c_d--project"></a>
### Nested Schema for `spec.argo_c_d.project`

Read-Only:

- `enabled` (Boolean)
- `metadata` (List of Object) (see [below for nested schema](#nestedobjatt--spec--argo_c_d--project--metadata))
- `roles` (List of Object) (see [below for nested schema](#nestedobjatt--spec--argo_c_d--project--roles))
- `source_repos` (List of String)

<a id="nestedobjatt--spec--argo_c_d--project--metadata"></a>
### Nested Schema for `spec.argo_c_d.project.source_repos`

Read-Only:

- `description` (String)
- `extra_annotations` (Map of String)
- `extra_labels` (Map of String)


<a id="nestedobjatt--spec--argo_c_d--project--roles"></a>
### Nested Schema for `spec.argo_c_d.project.source_repos`

Read-Only:

- `description` (String)
- `groups` (List of String)
- `name` (String)
- `rules` (List of Object) (see [below for nested schema](#nestedobjatt--spec--argo_c_d--project--source_repos--rules))

<a id="nestedobjatt--spec--argo_c_d--project--source_repos--rules"></a>
### Nested Schema for `spec.argo_c_d.project.source_repos.rules`

Read-Only:

- `action` (String)
- `application` (String)
- `permission` (Boolean)




<a id="nestedobjatt--spec--argo_c_d--sso"></a>
### Nested Schema for `spec.argo_c_d.sso`

Read-Only:

- `assigned_roles` (List of String)
- `enabled` (Boolean)
- `host` (String)



<a id="nestedobjatt--spec--members"></a>
### Nested Schema for `spec.members`

Read-Only:

- `cluster_role` (String)
- `group` (String)
- `kind` (String)
- `name` (String)


<a id="nestedobjatt--spec--namespace_pattern"></a>
### Nested Schema for `spec.namespace_pattern`

Read-Only:

- `space` (String)
- `virtual_cluster` (String)


<a id="nestedobjatt--spec--owner"></a>
### Nested Schema for `spec.owner`

Read-Only:

- `team` (String)
- `user` (String)


<a id="nestedobjatt--spec--quotas"></a>
### Nested Schema for `spec.quotas`

Read-Only:

- `project` (Map of String)
- `user` (Map of String)

