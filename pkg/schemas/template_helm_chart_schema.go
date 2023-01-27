//// Code generated by go-swagger; DO NOT EDIT.

package schemas

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	agentstoragev1 "github.com/loft-sh/agentapi/v2/pkg/apis/loft/storage/v1"
	"github.com/loft-sh/terraform-provider-loft/pkg/utils"
)

func StorageV1TemplateHelmChartSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"insecure_skip_tls_verify": {
			Type:        schema.TypeBool,
			Description: "If tls certificate checks for the chart download should be skipped",
			Optional:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Name is the chart name in the repository",
			Optional:    true,
		},
		"password": {
			Type:        schema.TypeString,
			Description: "The password that is required for this repository",
			Optional:    true,
		},
		"release_name": {
			Type:        schema.TypeString,
			Description: "ReleaseName is the preferred release name of the app",
			Optional:    true,
		},
		"release_namespace": {
			Type:        schema.TypeString,
			Description: "ReleaseNamespace is the preferred release namespace of the app",
			Optional:    true,
		},
		"repo_url": {
			Type:        schema.TypeString,
			Description: "RepoURL is the repo url where the chart can be found",
			Optional:    true,
		},
		"timeout": {
			Type:        schema.TypeString,
			Description: "Timeout is the time to wait for any individual Kubernetes operation (like Jobs for hooks) (default 5m0s)",
			Optional:    true,
		},
		"username": {
			Type:        schema.TypeString,
			Description: "The username that is required for this repository",
			Optional:    true,
		},
		"values": {
			Type:        schema.TypeString,
			Description: "Values are the values that should get passed to the chart",
			Optional:    true,
		},
		"version": {
			Type:        schema.TypeString,
			Description: "Version is the chart version in the repository",
			Optional:    true,
		},
		"wait": {
			Type:        schema.TypeBool,
			Description: "Wait determines if Loft should wait during deploy for the app to become ready",
			Optional:    true,
		},
	}
}

func CreateStorageV1TemplateHelmChart(data map[string]interface{}) *agentstoragev1.TemplateHelmChart {
	if !utils.HasKeys(data) {
		return nil
	}

	ret := &agentstoragev1.TemplateHelmChart{}

	if v, ok := data["insecure_skip_tls_verify"].(bool); ok {
		ret.InsecureSkipTlsVerify = v
	}

	if v, ok := data["name"].(string); ok && len(v) > 0 {
		ret.Name = v
	}

	if v, ok := data["password"].(string); ok && len(v) > 0 {
		ret.Password = v
	}

	if v, ok := data["release_name"].(string); ok && len(v) > 0 {
		ret.ReleaseName = v
	}

	if v, ok := data["release_namespace"].(string); ok && len(v) > 0 {
		ret.ReleaseNamespace = v
	}

	if v, ok := data["repo_url"].(string); ok && len(v) > 0 {
		ret.RepoURL = v
	}

	if v, ok := data["timeout"].(string); ok && len(v) > 0 {
		ret.Timeout = v
	}

	if v, ok := data["username"].(string); ok && len(v) > 0 {
		ret.Username = v
	}

	if v, ok := data["values"].(string); ok && len(v) > 0 {
		ret.Values = v
	}

	if v, ok := data["version"].(string); ok && len(v) > 0 {
		ret.Version = v
	}

	if v, ok := data["wait"].(bool); ok {
		ret.Wait = v
	}

	return ret
}

func ReadStorageV1TemplateHelmChart(obj *agentstoragev1.TemplateHelmChart) (interface{}, error) {
	if obj == nil {
		return nil, nil
	}

	values := map[string]interface{}{}
	values["insecure_skip_tls_verify"] = obj.InsecureSkipTlsVerify
	values["name"] = obj.Name
	values["password"] = obj.Password
	values["release_name"] = obj.ReleaseName
	values["release_namespace"] = obj.ReleaseNamespace
	values["repo_url"] = obj.RepoURL
	values["timeout"] = obj.Timeout
	values["username"] = obj.Username
	values["values"] = obj.Values
	values["version"] = obj.Version
	values["wait"] = obj.Wait
	return values, nil
}
