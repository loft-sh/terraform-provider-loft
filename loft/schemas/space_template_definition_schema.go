//// Code generated by go-swagger; DO NOT EDIT.

package schemas

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func StorageV1SpaceTemplateDefinitionSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"access": {
			Type:     schema.TypeList,
			MinItems: 1,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: StorageV1InstanceAccessSchema(),
			},
			Description: "The space access",
			Optional:    true,
		},
		"apps": {
			Type: schema.TypeList,
			Elem: &schema.Resource{
				Schema: StorageV1AppReferenceSchema(),
			},
			Description: "Apps specifies the apps that should get deployed by this template",
			Optional:    true,
		},
		"charts": {
			Type: schema.TypeList,
			Elem: &schema.Resource{
				Schema: StorageV1TemplateHelmChartSchema(),
			},
			Description: "Charts are helm charts that should get deployed",
			Optional:    true,
		},
		"metadata": {
			Type:     schema.TypeList,
			MinItems: 1,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: StorageV1TemplateMetadataSchema(),
			},
			Description: "The space metadata",
			Optional:    true,
		},
		"objects": {
			Type:        schema.TypeString,
			Description: "Objects are Kubernetes style yamls that should get deployed into the virtual cluster",
			Optional:    true,
		},
	}
}
