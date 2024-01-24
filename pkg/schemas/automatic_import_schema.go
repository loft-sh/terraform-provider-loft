//// Code generated by go-swagger; DO NOT EDIT.

package schemas

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	storagev1 "github.com/loft-sh/api/v3/pkg/apis/storage/v1"
	"github.com/loft-sh/terraform-provider-loft/pkg/utils"
)

func StorageV1AutomaticImportSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"virtual_clusters": {
			Type:     schema.TypeList,
			MinItems: 1,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: StorageV1AutomaticImportVirtualClustersSchema(),
			},
			Description: "VirtualClusters defines automatic virtual cluster import options.",
			Optional:    true,
			Computed:    true,
		},
	}
}

func CreateStorageV1AutomaticImport(data map[string]interface{}) *storagev1.AutomaticImport {
	if !utils.HasKeys(data) {
		return nil
	}

	ret := &storagev1.AutomaticImport{}

	if v, ok := data["virtual_clusters"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
		ret.VirtualClusters = *CreateStorageV1AutomaticImportVirtualClusters(v[0].(map[string]interface{}))
	}

	return ret
}

func ReadStorageV1AutomaticImport(obj *storagev1.AutomaticImport) (interface{}, error) {
	if obj == nil {
		return nil, nil
	}

	values := map[string]interface{}{}

	virtualClusters, err := ReadStorageV1AutomaticImportVirtualClusters(&obj.VirtualClusters)
	if err != nil {
		return nil, err
	}
	if virtualClusters != nil {
		values["virtual_clusters"] = []interface{}{virtualClusters}
	}

	return values, nil
}
