//// Code generated by go-swagger; DO NOT EDIT.

package schemas

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	storagev1 "github.com/loft-sh/api/v3/pkg/apis/storage/v1"
	"github.com/loft-sh/terraform-provider-loft/pkg/utils"
)

func StorageV1WorkloadVirtualClusterTemplateDefinitionSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"helm_release": {
			Type:     schema.TypeList,
			MinItems: 1,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: StorageV1VirtualClusterHelmReleaseSchema(),
			},
			Description: "HelmRelease is the helm release configuration for the virtual cluster.",
			Optional:    true,
		},
		"metadata": {
			Type:     schema.TypeList,
			MinItems: 1,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: StorageV1TemplateMetadataSchema(),
			},
			Description: "The virtual cluster metadata",
			Optional:    true,
		},
		"space_template": {
			Type:     schema.TypeList,
			MinItems: 1,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: StorageV1VirtualClusterSpaceTemplateDefinitionSchema(),
			},
			Description: "SpaceTemplate holds the space template",
			Optional:    true,
			Computed:    true,
		},
	}
}

func CreateStorageV1WorkloadVirtualClusterTemplateDefinition(data map[string]interface{}) *storagev1.WorkloadVirtualClusterTemplateDefinition {
	if !utils.HasKeys(data) {
		return nil
	}

	ret := &storagev1.WorkloadVirtualClusterTemplateDefinition{}

	if v, ok := data["helm_release"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
		ret.HelmRelease = CreateStorageV1VirtualClusterHelmRelease(v[0].(map[string]interface{}))
	}

	if v, ok := data["metadata"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
		ret.TemplateMetadata = *CreateStorageV1TemplateMetadata(v[0].(map[string]interface{}))
	}

	if v, ok := data["space_template"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
		ret.SpaceTemplate = *CreateStorageV1VirtualClusterSpaceTemplateDefinition(v[0].(map[string]interface{}))
	}

	return ret
}

func ReadStorageV1WorkloadVirtualClusterTemplateDefinition(obj *storagev1.WorkloadVirtualClusterTemplateDefinition) (interface{}, error) {
	if obj == nil {
		return nil, nil
	}

	values := map[string]interface{}{}
	helmRelease, err := ReadStorageV1VirtualClusterHelmRelease(obj.HelmRelease)
	if err != nil {
		return nil, err
	}
	if helmRelease != nil {
		values["helm_release"] = []interface{}{helmRelease}
	}

	metadata, err := ReadStorageV1TemplateMetadata(&obj.TemplateMetadata)
	if err != nil {
		return nil, err
	}
	if metadata != nil {
		values["metadata"] = []interface{}{metadata}
	}

	spaceTemplate, err := ReadStorageV1VirtualClusterSpaceTemplateDefinition(&obj.SpaceTemplate)
	if err != nil {
		return nil, err
	}
	if spaceTemplate != nil {
		values["space_template"] = []interface{}{spaceTemplate}
	}

	return values, nil
}
