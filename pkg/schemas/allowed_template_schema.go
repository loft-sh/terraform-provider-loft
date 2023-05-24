//// Code generated by go-swagger; DO NOT EDIT.

package schemas

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	storagev1 "github.com/loft-sh/api/v3/pkg/apis/storage/v1"
	"github.com/loft-sh/terraform-provider-loft/pkg/utils"
)

func StorageV1AllowedTemplateSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"group": {
			Type:        schema.TypeString,
			Description: "Group of the template that is allowed. Currently only supports storage.loft.sh",
			Optional:    true,
		},
		"is_default": {
			Type:        schema.TypeBool,
			Description: "IsDefault specifies if the template should be used as a default",
			Optional:    true,
		},
		"kind": {
			Type:        schema.TypeString,
			Description: "Kind of the template that is allowed. Currently only supports VirtualClusterTemplate & SpaceTemplate",
			Optional:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Name of the template",
			Optional:    true,
		},
	}
}

func CreateStorageV1AllowedTemplate(data map[string]interface{}) *storagev1.AllowedTemplate {
	if !utils.HasKeys(data) {
		return nil
	}

	ret := &storagev1.AllowedTemplate{}
	if v, ok := data["group"].(string); ok && len(v) > 0 {
		ret.Group = v
	}

	if v, ok := data["is_default"].(bool); ok {
		ret.IsDefault = v
	}

	if v, ok := data["kind"].(string); ok && len(v) > 0 {
		ret.Kind = v
	}

	if v, ok := data["name"].(string); ok && len(v) > 0 {
		ret.Name = v
	}

	return ret
}

func ReadStorageV1AllowedTemplate(obj *storagev1.AllowedTemplate) (interface{}, error) {
	if obj == nil {
		return nil, nil
	}

	values := map[string]interface{}{}
	values["group"] = obj.Group

	values["is_default"] = obj.IsDefault

	values["kind"] = obj.Kind

	values["name"] = obj.Name

	return values, nil
}
