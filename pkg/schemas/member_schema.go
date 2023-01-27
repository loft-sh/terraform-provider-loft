//// Code generated by go-swagger; DO NOT EDIT.

package schemas

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	storagev1 "github.com/loft-sh/api/v2/pkg/apis/storage/v1"
	"github.com/loft-sh/terraform-provider-loft/pkg/utils"
)

func StorageV1MemberSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cluster_role": {
			Type:        schema.TypeString,
			Description: "ClusterRole is the assigned role for the above member",
			Optional:    true,
		},
		"group": {
			Type:        schema.TypeString,
			Description: "Group of the member. Currently only supports storage.loft.sh",
			Optional:    true,
		},
		"kind": {
			Type:        schema.TypeString,
			Description: "Kind is the kind of the member. Currently either User or Team",
			Optional:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Name of the member",
			Optional:    true,
		},
	}
}

func CreateStorageV1Member(data map[string]interface{}) *storagev1.Member {
	if !utils.HasKeys(data) {
		return nil
	}

	ret := &storagev1.Member{}
	if v, ok := data["cluster_role"].(string); ok && len(v) > 0 {
		ret.ClusterRole = v
	}

	if v, ok := data["group"].(string); ok && len(v) > 0 {
		ret.Group = v
	}

	if v, ok := data["kind"].(string); ok && len(v) > 0 {
		ret.Kind = v
	}

	if v, ok := data["name"].(string); ok && len(v) > 0 {
		ret.Name = v
	}

	return ret
}

func ReadStorageV1Member(obj *storagev1.Member) (interface{}, error) {
	if obj == nil {
		return nil, nil
	}

	values := map[string]interface{}{}
	values["cluster_role"] = obj.ClusterRole

	values["group"] = obj.Group

	values["kind"] = obj.Kind

	values["name"] = obj.Name

	return values, nil
}
