//// Code generated by go-swagger; DO NOT EDIT.

package schemas

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/loft-sh/terraform-provider-loft/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"
)

func IoK8sAPICoreV1SecretKeySelectorSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"key": {
			Type:        schema.TypeString,
			Description: "The key of the secret to select from.  Must be a valid secret key.",
			Required:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names",
			Optional:    true,
		},
		"optional": {
			Type:        schema.TypeBool,
			Description: "Specify whether the Secret or its key must be defined",
			Optional:    true,
		},
	}
}

func CreateIoK8sAPICoreV1SecretKeySelector(data map[string]interface{}) *corev1.SecretKeySelector {
	if !utils.HasKeys(data) {
		return nil
	}

	ret := &corev1.SecretKeySelector{}
	if v, ok := data["key"].(string); ok && len(v) > 0 {
		ret.Key = v
	}

	if v, ok := data["name"].(string); ok && len(v) > 0 {
		ret.Name = v
	}

	if v, ok := data["optional"].(bool); ok {
		ret.Optional = ptr.To(v)
	}

	return ret
}

func ReadIoK8sAPICoreV1SecretKeySelector(obj *corev1.SecretKeySelector) (interface{}, error) {
	if obj == nil {
		return nil, nil
	}

	values := map[string]interface{}{}
	values["key"] = obj.Key

	values["name"] = obj.Name

	values["optional"] = obj.Optional

	return values, nil
}
