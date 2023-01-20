//// Code generated by go-swagger; DO NOT EDIT.

package schemas

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	agentstoragev1 "github.com/loft-sh/agentapi/v2/pkg/apis/loft/storage/v1"
	"github.com/loft-sh/terraform-provider-loft/pkg/utils"
)

func StorageV1InstanceAccessRuleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cluster_role": {
			Type:        schema.TypeString,
			Description: "ClusterRole is the cluster role that should be assigned to the",
			Optional:    true,
		},
		"teams": {
			Type: schema.TypeList,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: "Teams that this rule matches.",
			Optional:    true,
		},
		"users": {
			Type: schema.TypeList,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: "Users this rule matches. * means all users.",
			Optional:    true,
		},
	}
}

func CreateStorageV1InstanceAccessRule(in []interface{}) *agentstoragev1.InstanceAccessRule {
	if !utils.HasValue(in) {
		return nil
	}

	ret := &agentstoragev1.InstanceAccessRule{}

	data := in[0].(map[string]interface{})
	if v, ok := data["cluster_role"].(string); ok && len(v) > 0 {
		ret.ClusterRole = v
	}

	var teamsItems []string
	for _, v := range data["teams"].([]string) {
		teamsItems = append(teamsItems, v)
	}
	ret.Teams = teamsItems

	var usersItems []string
	for _, v := range data["users"].([]string) {
		usersItems = append(usersItems, v)
	}
	ret.Users = usersItems

	return ret
}

func ReadStorageV1InstanceAccessRule(obj *agentstoragev1.InstanceAccessRule) (interface{}, error) {
	values := map[string]interface{}{}
	values["cluster_role"] = obj.ClusterRole
	var teamsItems []interface{}
	for _, v := range obj.Teams {
		teamsItems = append(teamsItems, v)
	}
	values["teams"] = teamsItems
	var usersItems []interface{}
	for _, v := range obj.Users {
		usersItems = append(usersItems, v)
	}
	values["users"] = usersItems
	return values, nil
}
