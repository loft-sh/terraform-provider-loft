//// Code generated by go-swagger; DO NOT EDIT.

package schemas

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	managementv1 "github.com/loft-sh/api/v3/pkg/apis/management/v1"
	storagev1 "github.com/loft-sh/api/v3/pkg/apis/storage/v1"
	"github.com/loft-sh/terraform-provider-loft/pkg/utils"
)

func ManagementV1VirtualClusterTemplateSpecSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"access": {
			Type: schema.TypeList,
			Elem: &schema.Resource{
				Schema: StorageV1AccessSchema(),
			},
			Description: "Access holds the access rights for users and teams",
			Optional:    true,
		},
		"description": {
			Type:        schema.TypeString,
			Description: "Description describes the virtual cluster template",
			Optional:    true,
		},
		"display_name": {
			Type:        schema.TypeString,
			Description: "DisplayName is the name that is shown in the UI",
			Optional:    true,
		},
		"owner": {
			Type:     schema.TypeList,
			MinItems: 1,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: StorageV1UserOrTeamSchema(),
			},
			Description: "Owner holds the owner of this object",
			Optional:    true,
		},
		"parameters": {
			Type: schema.TypeList,
			Elem: &schema.Resource{
				Schema: StorageV1AppParameterSchema(),
			},
			Description: "Parameters define additional app parameters that will set helm values",
			Optional:    true,
		},
		"space_template_ref": {
			Type:     schema.TypeList,
			MinItems: 1,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: StorageV1VirtualClusterTemplateSpaceTemplateRefSchema(),
			},
			Description: "DEPRECATED: SpaceTemplate to use to create the virtual cluster space if it does not exist",
			Optional:    true,
		},
		"template": {
			Type:     schema.TypeList,
			MinItems: 1,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: StorageV1VirtualClusterTemplateDefinitionSchema(),
			},
			Description: "Template holds the virtual cluster template",
			Optional:    true,
		},
		"versions": {
			Type: schema.TypeList,
			Elem: &schema.Resource{
				Schema: StorageV1VirtualClusterTemplateVersionSchema(),
			},
			Description: "Versions are different versions of the template that can be referenced as well",
			Optional:    true,
		},
	}
}

func CreateManagementV1VirtualClusterTemplateSpec(data map[string]interface{}) *managementv1.VirtualClusterTemplateSpec {
	ret := storagev1.VirtualClusterTemplateSpec{}

	if utils.HasKeys(data) {

		var accessItems []storagev1.Access
		for _, v := range data["access"].([]interface{}) {
			if v == nil {
				continue
			}
			if item := CreateStorageV1Access(v.(map[string]interface{})); item != nil {
				accessItems = append(accessItems, *item)
			}
		}
		ret.Access = accessItems

		if v, ok := data["description"].(string); ok && len(v) > 0 {
			ret.Description = v
		}

		if v, ok := data["display_name"].(string); ok && len(v) > 0 {
			ret.DisplayName = v
		}

		if v, ok := data["owner"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
			ret.Owner = CreateStorageV1UserOrTeam(v[0].(map[string]interface{}))
		}

		var parametersItems []storagev1.AppParameter
		for _, v := range data["parameters"].([]interface{}) {
			if v == nil {
				continue
			}
			if item := CreateStorageV1AppParameter(v.(map[string]interface{})); item != nil {
				parametersItems = append(parametersItems, *item)
			}
		}
		ret.Parameters = parametersItems

		if v, ok := data["space_template_ref"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
			ret.SpaceTemplateRef = CreateStorageV1VirtualClusterTemplateSpaceTemplateRef(v[0].(map[string]interface{}))
		}

		if v, ok := data["template"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
			ret.Template = *CreateStorageV1VirtualClusterTemplateDefinition(v[0].(map[string]interface{}))
		}

		var versionsItems []storagev1.VirtualClusterTemplateVersion
		for _, v := range data["versions"].([]interface{}) {
			if v == nil {
				continue
			}
			if item := CreateStorageV1VirtualClusterTemplateVersion(v.(map[string]interface{})); item != nil {
				versionsItems = append(versionsItems, *item)
			}
		}
		ret.Versions = versionsItems

	}

	return &managementv1.VirtualClusterTemplateSpec{
		VirtualClusterTemplateSpec: ret,
	}
}

func ReadManagementV1VirtualClusterTemplateSpec(obj *managementv1.VirtualClusterTemplateSpec) (interface{}, error) {
	if obj == nil {
		return nil, nil
	}

	values := map[string]interface{}{}
	var accessItems []interface{}
	for _, v := range obj.Access {
		item, err := ReadStorageV1Access(&v)
		if err != nil {
			return nil, err
		}
		accessItems = append(accessItems, item)
	}
	values["access"] = accessItems

	values["description"] = obj.Description

	values["display_name"] = obj.DisplayName

	owner, err := ReadStorageV1UserOrTeam(obj.Owner)
	if err != nil {
		return nil, err
	}
	if owner != nil {
		values["owner"] = []interface{}{owner}
	}

	var parametersItems []interface{}
	for _, v := range obj.Parameters {
		item, err := ReadStorageV1AppParameter(&v)
		if err != nil {
			return nil, err
		}
		parametersItems = append(parametersItems, item)
	}
	values["parameters"] = parametersItems

	spaceTemplateRef, err := ReadStorageV1VirtualClusterTemplateSpaceTemplateRef(obj.SpaceTemplateRef)
	if err != nil {
		return nil, err
	}
	values["space_template_ref"] = spaceTemplateRef

	template, err := ReadStorageV1VirtualClusterTemplateDefinition(&obj.Template)
	if err != nil {
		return nil, err
	}
	if template != nil {
		values["template"] = []interface{}{template}
	}

	var versionsItems []interface{}
	for _, v := range obj.Versions {
		item, err := ReadStorageV1VirtualClusterTemplateVersion(&v)
		if err != nil {
			return nil, err
		}
		versionsItems = append(versionsItems, item)
	}
	values["versions"] = versionsItems

	return values, nil
}