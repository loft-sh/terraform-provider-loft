//// Code generated by go-swagger; DO NOT EDIT.

package schemas

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	storagev1 "github.com/loft-sh/api/v2/pkg/apis/storage/v1"
	"github.com/loft-sh/terraform-provider-loft/pkg/utils"
)

func StorageV1VirtualClusterClusterRefSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cluster": {
			Type:        schema.TypeString,
			Description: "Cluster is the connected cluster the space will be created in",
			Optional:    true,
		},
		"namespace": {
			Type:        schema.TypeString,
			Description: "Namespace is the namespace inside the connected cluster holding the space",
			Optional:    true,
		},
		"virtual_cluster": {
			Type:        schema.TypeString,
			Description: "VirtualCluster is the name of the virtual cluster inside the namespace",
			Optional:    true,
		},
	}
}

func CreateStorageV1VirtualClusterClusterRef(in []interface{}) *storagev1.VirtualClusterClusterRef {
	if !utils.HasValue(in) {
		return nil
	}

	ret := &storagev1.VirtualClusterClusterRef{}

	data := in[0].(map[string]interface{})
	if v, ok := data["cluster"].(string); ok && len(v) > 0 {
		ret.Cluster = v
	}

	if v, ok := data["namespace"].(string); ok && len(v) > 0 {
		ret.Namespace = v
	}

	if v, ok := data["virtual_cluster"].(string); ok && len(v) > 0 {
		ret.VirtualCluster = v
	}

	return ret
}

func ReadStorageV1VirtualClusterClusterRef(obj *storagev1.VirtualClusterClusterRef) (interface{}, error) {
	values := map[string]interface{}{}
	values["cluster"] = obj.Cluster
	values["namespace"] = obj.Namespace
	values["virtual_cluster"] = obj.VirtualCluster
	return values, nil
}
