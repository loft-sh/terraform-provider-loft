//// Code generated by go-swagger; DO NOT EDIT.

package schemas

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	agentstoragev1 "github.com/loft-sh/agentapi/v2/pkg/apis/loft/storage/v1"
	managementv1 "github.com/loft-sh/api/v2/pkg/apis/management/v1"
	storagev1 "github.com/loft-sh/api/v2/pkg/apis/storage/v1"
	"github.com/loft-sh/terraform-provider-loft/pkg/utils"
)

func ManagementV1VirtualClusterInstanceSpecSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"access": {
			Type: schema.TypeList,
			Elem: &schema.Resource{
				Schema: StorageV1AccessSchema(),
			},
			Description: "Access to the virtual cluster object itself",
			Optional:    true,
			Computed:    true,
		},
		"cluster_ref": {
			Type:     schema.TypeList,
			MinItems: 1,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: StorageV1VirtualClusterClusterRefSchema(),
			},
			Description: "ClusterRef is the reference to the connected cluster holding this virtual cluster",
			Optional:    true,
			Computed:    true,
		},
		"description": {
			Type:        schema.TypeString,
			Description: "Description describes a virtual cluster instance",
			Optional:    true,
		},
		"display_name": {
			Type:        schema.TypeString,
			Description: "DisplayName is the name that should be displayed in the UI",
			Optional:    true,
		},
		"extra_access_rules": {
			Type: schema.TypeList,
			Elem: &schema.Resource{
				Schema: StorageV1InstanceAccessRuleSchema(),
			},
			Description: "ExtraAccessRules defines extra rules which users and teams should have which access to the virtual cluster.",
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
			Type:        schema.TypeString,
			Description: "Parameters are values to pass to the template",
			Optional:    true,
		},
		"template": {
			Type:     schema.TypeList,
			MinItems: 1,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: StorageV1VirtualClusterTemplateDefinitionSchema(),
			},
			Description: "Template is the inline template to use for virtual cluster creation. This is mutually exclusive with templateRef.",
			Optional:    true,
		},
		"template_ref": {
			Type:     schema.TypeList,
			MinItems: 1,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: StorageV1TemplateRefSchema(),
			},
			Description: "TemplateRef holds the virtual cluster template reference",
			Optional:    true,
		},
	}
}

func CreateManagementV1VirtualClusterInstanceSpec(data map[string]interface{}) *managementv1.VirtualClusterInstanceSpec {
	ret := storagev1.VirtualClusterInstanceSpec{}

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

		if v, ok := data["cluster_ref"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
			ret.ClusterRef = *CreateStorageV1VirtualClusterClusterRef(v[0].(map[string]interface{}))
		}

		if v, ok := data["description"].(string); ok && len(v) > 0 {
			ret.Description = v
		}

		if v, ok := data["display_name"].(string); ok && len(v) > 0 {
			ret.DisplayName = v
		}

		var extraAccessRulesItems []agentstoragev1.InstanceAccessRule
		for _, v := range data["extra_access_rules"].([]interface{}) {
			if v == nil {
				continue
			}
			if item := CreateStorageV1InstanceAccessRule(v.(map[string]interface{})); item != nil {
				extraAccessRulesItems = append(extraAccessRulesItems, *item)
			}
		}
		ret.ExtraAccessRules = extraAccessRulesItems

		if v, ok := data["owner"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
			ret.Owner = CreateStorageV1UserOrTeam(v[0].(map[string]interface{}))
		}

		if v, ok := data["parameters"].(string); ok && len(v) > 0 {
			ret.Parameters = v
		}

		if v, ok := data["template"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
			ret.Template = CreateStorageV1VirtualClusterTemplateDefinition(v[0].(map[string]interface{}))
		}

		if v, ok := data["template_ref"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
			ret.TemplateRef = CreateStorageV1TemplateRef(v[0].(map[string]interface{}))
		}

	}

	return &managementv1.VirtualClusterInstanceSpec{
		VirtualClusterInstanceSpec: ret,
	}
}

func ReadManagementV1VirtualClusterInstanceSpec(obj *managementv1.VirtualClusterInstanceSpec) (interface{}, error) {
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

	clusterRef, err := ReadStorageV1VirtualClusterClusterRef(&obj.ClusterRef)
	if err != nil {
		return nil, err
	}
	if clusterRef != nil {
		values["cluster_ref"] = []interface{}{clusterRef}
	}

	values["description"] = obj.Description

	values["display_name"] = obj.DisplayName

	var extraAccessRulesItems []interface{}
	for _, v := range obj.ExtraAccessRules {
		item, err := ReadStorageV1InstanceAccessRule(&v)
		if err != nil {
			return nil, err
		}
		extraAccessRulesItems = append(extraAccessRulesItems, item)
	}
	values["extra_access_rules"] = extraAccessRulesItems

	owner, err := ReadStorageV1UserOrTeam(obj.Owner)
	if err != nil {
		return nil, err
	}
	if owner != nil {
		values["owner"] = []interface{}{owner}
	}

	values["parameters"] = obj.Parameters

	template, err := ReadStorageV1VirtualClusterTemplateDefinition(obj.Template)
	if err != nil {
		return nil, err
	}
	if template != nil {
		values["template"] = []interface{}{template}
	}

	templateRef, err := ReadStorageV1TemplateRef(obj.TemplateRef)
	if err != nil {
		return nil, err
	}
	if templateRef != nil {
		values["template_ref"] = []interface{}{templateRef}
	}

	return values, nil
}
