//// Code generated by go-swagger; DO NOT EDIT.

package schemas

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	storagev1 "github.com/loft-sh/api/v3/pkg/apis/storage/v1"
	"github.com/loft-sh/terraform-provider-loft/pkg/utils"
)

func StorageV1ArgoSSOSpecSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"assigned_roles": {
			Type: schema.TypeList,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: "AssignedRoles is a list of roles to assign for users who authenticate via Loft -- by default this will be the `read-only` role. If any roles are provided this will override the default setting.",
			Optional:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "Enabled indicates if the ArgoCD SSO Integration is enabled for this project. Enabling this will cause Loft to configure SSO authentication via Loft in ArgoCD. If Projects are *not* enabled, all users associated with this Project will be assigned either the 'read-only' (default) role, *or* the roles set under the AssignedRoles field.",
			Optional:    true,
		},
		"host": {
			Type:        schema.TypeString,
			Description: "Host defines the ArgoCD host address that will be used for OIDC authentication between loft and ArgoCD. If not specified OIDC integration will be skipped, but vclusters/spaces will still be synced to ArgoCD.",
			Optional:    true,
		},
	}
}

func CreateStorageV1ArgoSSOSpec(data map[string]interface{}) *storagev1.ArgoSSOSpec {
	if !utils.HasKeys(data) {
		return nil
	}

	ret := &storagev1.ArgoSSOSpec{}
	var assignedRolesItems []string
	for _, v := range data["assigned_roles"].([]interface{}) {
		assignedRolesItems = append(assignedRolesItems, v.(string))
	}
	ret.AssignedRoles = assignedRolesItems

	if v, ok := data["enabled"].(bool); ok {
		ret.Enabled = v
	}

	if v, ok := data["host"].(string); ok && len(v) > 0 {
		ret.Host = v
	}

	return ret
}

func ReadStorageV1ArgoSSOSpec(obj *storagev1.ArgoSSOSpec) (interface{}, error) {
	if obj == nil {
		return nil, nil
	}

	values := map[string]interface{}{}
	var assignedRolesItems []interface{}
	for _, v := range obj.AssignedRoles {
		assignedRolesItems = append(assignedRolesItems, v)
	}
	values["assigned_roles"] = assignedRolesItems

	values["enabled"] = obj.Enabled

	values["host"] = obj.Host

	return values, nil
}
