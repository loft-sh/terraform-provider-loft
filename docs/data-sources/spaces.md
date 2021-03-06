---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "loft_spaces Data Source - terraform-provider-loft"
subcategory: ""
description: |-
  Provides details for loft_spaces Data Source
---

# loft_spaces (Data Source)

The `loft_spaces` data source provides information about all Loft spaces in the given `cluster`.

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

# Import data for all spaces
data "loft_spaces" "all" {
  cluster = "loft-cluster"
}

# Output all spaces
output "spaces" {
  value = data.loft_spaces.all.spaces.*.name
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cluster` (String) The cluster to list spaces from.

### Read-Only

- `id` (String) The ID of this resource.
- `spaces` (List of Object) All spaces (see [below for nested schema](#nestedatt--spaces))

<a id="nestedatt--spaces"></a>
### Nested Schema for `spaces`

Read-Only:

- `annotations` (Map of String)
- `cluster` (String)
- `delete_after` (String)
- `generate_name` (String)
- `id` (String)
- `labels` (Map of String)
- `name` (String)
- `objects` (String)
- `sleep_after` (String)
- `sleep_schedule` (String)
- `space_constraints` (String)
- `team` (String)
- `user` (String)
- `wakeup_schedule` (String)


