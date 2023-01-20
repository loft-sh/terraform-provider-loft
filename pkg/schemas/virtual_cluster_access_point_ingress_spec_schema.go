//// Code generated by go-swagger; DO NOT EDIT.

package schemas

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	agentstoragev1 "github.com/loft-sh/agentapi/v2/pkg/apis/loft/storage/v1"
	"github.com/loft-sh/terraform-provider-loft/pkg/utils"
)

func StorageV1VirtualClusterAccessPointIngressSpecSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "Enabled defines if the virtual cluster access point (via ingress) is enabled or not; requires the connected cluster to have the `loft.sh/ingress-suffix` annotation set to define the domain name suffix used for the ingress.",
			Optional:    true,
		},
	}
}

func CreateStorageV1VirtualClusterAccessPointIngressSpec(in []interface{}) *agentstoragev1.VirtualClusterAccessPointIngressSpec {
	if !utils.HasValue(in) {
		return nil
	}

	ret := &agentstoragev1.VirtualClusterAccessPointIngressSpec{}

	data := in[0].(map[string]interface{})
	if v, ok := data["enabled"].(bool); ok {
		ret.Enabled = v
	}

	return ret
}

func ReadStorageV1VirtualClusterAccessPointIngressSpec(obj *agentstoragev1.VirtualClusterAccessPointIngressSpec) (interface{}, error) {
	values := map[string]interface{}{}
	values["enabled"] = obj.Enabled
	return values, nil
}
