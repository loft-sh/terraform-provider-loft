//// Code generated by go-swagger; DO NOT EDIT.

package schemas

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	agentstoragev1 "github.com/loft-sh/agentapi/v2/pkg/apis/loft/storage/v1"
	"github.com/loft-sh/terraform-provider-loft/pkg/utils"
)

func StorageV1VirtualClusterHelmReleaseSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"chart": {
			Type:     schema.TypeList,
			MinItems: 1,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: StorageV1VirtualClusterHelmChartSchema(),
			},
			Description: "infos about what chart to deploy",
			Optional:    true,
		},
		"values": {
			Type:        schema.TypeString,
			Description: "the values for the given chart",
			Optional:    true,
		},
	}
}

func CreateStorageV1VirtualClusterHelmRelease(data map[string]interface{}) *agentstoragev1.VirtualClusterHelmRelease {
	if !utils.HasKeys(data) {
		return nil
	}

	ret := &agentstoragev1.VirtualClusterHelmRelease{}

	if value := CreateStorageV1VirtualClusterHelmChart(data["chart"].(map[string]interface{})); value != nil {
		ret.Chart = *value
	}

	if v, ok := data["values"].(string); ok && len(v) > 0 {
		ret.Values = v
	}

	return ret
}

func ReadStorageV1VirtualClusterHelmRelease(obj *agentstoragev1.VirtualClusterHelmRelease) (interface{}, error) {
	if obj == nil {
		return nil, nil
	}

	values := map[string]interface{}{}

	chart, err := ReadStorageV1VirtualClusterHelmChart(&obj.Chart)
	if err != nil {
		return nil, err
	}
	values["chart"] = []interface{}{chart}

	values["values"] = obj.Values

	return values, nil
}
