//// Code generated by go-swagger; DO NOT EDIT.

package schemas

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	storagev1 "github.com/loft-sh/api/v3/pkg/apis/storage/v1"
	"github.com/loft-sh/terraform-provider-loft/pkg/utils"
)

func StorageV1ArgoProjectSpecMetadataSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Type:        schema.TypeString,
			Description: "Description to add to the ArgoCD project.",
			Optional:    true,
		},
		"extra_annotations": {
			Type: schema.TypeMap,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: "ExtraAnnotations are optional annotations that can be attached to the project in ArgoCD.",
			Optional:    true,
		},
		"extra_labels": {
			Type: schema.TypeMap,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: "ExtraLabels are optional labels that can be attached to the project in ArgoCD.",
			Optional:    true,
		},
	}
}

func CreateStorageV1ArgoProjectSpecMetadata(data map[string]interface{}) *storagev1.ArgoProjectSpecMetadata {
	if !utils.HasKeys(data) {
		return nil
	}

	ret := &storagev1.ArgoProjectSpecMetadata{}
	if v, ok := data["description"].(string); ok && len(v) > 0 {
		ret.Description = v
	}

	ret.ExtraAnnotations = utils.AttributesToMap(data["extra_annotations"].(map[string]interface{}))

	ret.ExtraLabels = utils.AttributesToMap(data["extra_labels"].(map[string]interface{}))

	return ret
}

func ReadStorageV1ArgoProjectSpecMetadata(obj *storagev1.ArgoProjectSpecMetadata) (interface{}, error) {
	if obj == nil {
		return nil, nil
	}

	values := map[string]interface{}{}
	values["description"] = obj.Description

	values["extra_annotations"] = obj.ExtraAnnotations

	values["extra_labels"] = obj.ExtraLabels

	return values, nil
}
