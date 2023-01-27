//// Code generated by go-swagger; DO NOT EDIT.

package schemas

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	managementv1 "github.com/loft-sh/api/v2/pkg/apis/management/v1"
	storagev1 "github.com/loft-sh/api/v2/pkg/apis/storage/v1"
	"github.com/loft-sh/terraform-provider-loft/pkg/utils"
)

func ManagementV1ProjectSpecSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"access": {
			Type: schema.TypeList,
			Elem: &schema.Resource{
				Schema: StorageV1AccessSchema(),
			},
			Description: "Access holds the access rights for users and teams",
			Optional:    true,
			Computed:    true,
		},
		"allowed_clusters": {
			Type: schema.TypeList,
			Elem: &schema.Resource{
				Schema: StorageV1AllowedClusterSchema(),
			},
			Description: "AllowedClusters are target clusters that are allowed to target with environments.",
			Optional:    true,
		},
		"allowed_templates": {
			Type: schema.TypeList,
			Elem: &schema.Resource{
				Schema: StorageV1AllowedTemplateSchema(),
			},
			Description: "AllowedTemplates are the templates that are allowed to use in this project.",
			Optional:    true,
		},
		"argo_c_d": {
			Type:     schema.TypeList,
			MinItems: 1,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: StorageV1ArgoIntegrationSpecSchema(),
			},
			Description: "ArgoIntegration holds information about ArgoCD Integration",
			Optional:    true,
		},
		"description": {
			Type:        schema.TypeString,
			Description: "Description describes an app",
			Optional:    true,
		},
		"display_name": {
			Type:        schema.TypeString,
			Description: "DisplayName is the name that should be displayed in the UI",
			Optional:    true,
		},
		"members": {
			Type: schema.TypeList,
			Elem: &schema.Resource{
				Schema: StorageV1MemberSchema(),
			},
			Description: "Members are the users and teams that are part of this project",
			Optional:    true,
		},
		"namespace_pattern": {
			Type:     schema.TypeList,
			MinItems: 1,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: StorageV1NamespacePatternSchema(),
			},
			Description: "NamespacePattern specifies template patterns to use for creating each space or virtual cluster's namespace",
			Optional:    true,
			Computed:    true,
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
		"quotas": {
			Type:     schema.TypeList,
			MinItems: 1,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: StorageV1QuotasSchema(),
			},
			Description: "Quotas define the quotas inside the project",
			Optional:    true,
			Computed:    true,
		},
	}
}

func CreateManagementV1ProjectSpec(data map[string]interface{}) *managementv1.ProjectSpec {
	ret := storagev1.ProjectSpec{}

	if utils.HasKeys(data) {
		var accessItems []storagev1.Access
		for _, v := range data["access"].([]interface{}) {
			if item := CreateStorageV1Access(v.(map[string]interface{})); item != nil {
				accessItems = append(accessItems, *item)
			}
		}
		ret.Access = accessItems

		var allowedClustersItems []storagev1.AllowedCluster
		for _, v := range data["allowed_clusters"].([]interface{}) {
			if v == nil {
				continue
			}

			if item := CreateStorageV1AllowedCluster(v.(map[string]interface{})); item != nil {
				allowedClustersItems = append(allowedClustersItems, *item)
			}
		}
		ret.AllowedClusters = allowedClustersItems

		var allowedTemplatesItems []storagev1.AllowedTemplate
		for _, v := range data["allowed_templates"].([]interface{}) {
			if item := CreateStorageV1AllowedTemplate(v.(map[string]interface{})); item != nil {
				allowedTemplatesItems = append(allowedTemplatesItems, *item)
			}
		}
		ret.AllowedTemplates = allowedTemplatesItems

		if v, ok := data["argo_c_d"].([]interface{}); ok && len(v) > 0 {
			ret.ArgoIntegration = CreateStorageV1ArgoIntegrationSpec(v[0].(map[string]interface{}))
		}

		if v, ok := data["description"].(string); ok && len(v) > 0 {
			ret.Description = v
		}

		if v, ok := data["display_name"].(string); ok && len(v) > 0 {
			ret.DisplayName = v
		}

		var membersItems []storagev1.Member
		for _, v := range data["members"].([]interface{}) {
			if item := CreateStorageV1Member(v.(map[string]interface{})); item != nil {
				membersItems = append(membersItems, *item)
			}
		}
		ret.Members = membersItems

		if v, ok := data["namespace_pattern"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
			ret.NamespacePattern = CreateStorageV1NamespacePattern(v[0].(map[string]interface{}))
		}

		if v, ok := data["owner"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
			ret.Owner = CreateStorageV1UserOrTeam(v[0].(map[string]interface{}))
		}

		if v, ok := data["quotas"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
			ret.Quotas = *CreateStorageV1Quotas(v[0].(map[string]interface{}))
		}
	}

	return &managementv1.ProjectSpec{
		ProjectSpec: ret,
	}
}

func ReadManagementV1ProjectSpec(obj *managementv1.ProjectSpec) (interface{}, error) {
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
	var allowedClustersItems []interface{}
	for _, v := range obj.AllowedClusters {
		item, err := ReadStorageV1AllowedCluster(&v)
		if err != nil {
			return nil, err
		}
		allowedClustersItems = append(allowedClustersItems, item)
	}
	values["allowed_clusters"] = allowedClustersItems
	var allowedTemplatesItems []interface{}
	for _, v := range obj.AllowedTemplates {
		item, err := ReadStorageV1AllowedTemplate(&v)
		if err != nil {
			return nil, err
		}
		allowedTemplatesItems = append(allowedTemplatesItems, item)
	}
	values["allowed_templates"] = allowedTemplatesItems

	argoCD, err := ReadStorageV1ArgoIntegrationSpec(obj.ArgoIntegration)
	if err != nil {
		return nil, err
	}
	values["argo_c_d"] = []interface{}{argoCD}
	values["description"] = obj.Description
	values["display_name"] = obj.DisplayName
	var membersItems []interface{}
	for _, v := range obj.Members {
		item, err := ReadStorageV1Member(&v)
		if err != nil {
			return nil, err
		}
		membersItems = append(membersItems, item)
	}
	values["members"] = membersItems

	namespacePattern, err := ReadStorageV1NamespacePattern(obj.NamespacePattern)
	if err != nil {
		return nil, err
	}
	values["namespace_pattern"] = []interface{}{namespacePattern}

	owner, err := ReadStorageV1UserOrTeam(obj.Owner)
	if err != nil {
		return nil, err
	}
	values["owner"] = []interface{}{owner}

	quotas, err := ReadStorageV1Quotas(&obj.Quotas)
	if err != nil {
		return nil, err
	}
	values["quotas"] = []interface{}{quotas}
	return values, nil
}
