//// Code generated by go-swagger; DO NOT EDIT.

package schemas

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	agentstoragev1 "github.com/loft-sh/agentapi/v3/pkg/apis/loft/storage/v1"
	"github.com/loft-sh/terraform-provider-loft/pkg/utils"
)

func StorageV1VirtualClusterProSpecSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "Enabled defines if the virtual cluster is a pro cluster or not",
			Optional:    true,
		},
	}
}

func CreateStorageV1VirtualClusterProSpec(data map[string]interface{}) *agentstoragev1.VirtualClusterProSpec {
	if !utils.HasKeys(data) {
		return nil
	}

	ret := &agentstoragev1.VirtualClusterProSpec{}
	if v, ok := data["enabled"].(bool); ok {
		ret.Enabled = v
	}

	return ret
}

func ReadStorageV1VirtualClusterProSpec(obj *agentstoragev1.VirtualClusterProSpec) (interface{}, error) {
	if obj == nil {
		return nil, nil
	}

	values := map[string]interface{}{}
	values["enabled"] = obj.Enabled

	return values, nil
}