package utils

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

func ReadId(metadata metav1.ObjectMeta) string {
	if metadata.Namespace != "" {
		return metadata.Namespace + "/" + metadata.Name
	}

	return metadata.Name
}

func ParseID(id string) (string, string) {
	tokens := strings.Split(id, "/")
	if len(tokens) == 2 {
		return tokens[0], tokens[1]
	}

	return "", ""
}

func MetadataSchema(objectName string, generateName bool, clusterScope bool) *schema.Schema {
	fields := metadataFields(objectName)

	if generateName {
		fields["generate_name"] = &schema.Schema{
			Type:          schema.TypeString,
			Description:   "Prefix, used by the server, to generate a unique name ONLY IF the `name` field has not been provided. This value will also be combined with a unique suffix. Read more: https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#idempotency",
			Optional:      true,
			ForceNew:      true,
			Computed:      true,
			ConflictsWith: []string{"metadata.0.name"},
			AtLeastOneOf:  []string{"metadata.0.name", "metadata.0.generate_name"},
		}
		fields["name"].ConflictsWith = []string{"metadata.0.generate_name"}
	}

	if !clusterScope {
		fields["namespace"] = &schema.Schema{
			Type:        schema.TypeString,
			Description: "Namespace defines the space within which each name must be unique. An empty namespace is equivalent to the \"default\" namespace, but \"default\" is the canonical representation. Not all objects are required to be scoped to a namespace - the value of this field for those objects will be empty.\n\nMust be a DNS_LABEL. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/namespaces",
			Required:    true,
		}
	}

	return &schema.Schema{
		Type:        schema.TypeList,
		Description: fmt.Sprintf("Standard %s's metadata. More info: https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#metadata", objectName),
		Required:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: fields,
		},
	}
}

func ReadMetadata(metadata metav1.ObjectMeta) (interface{}, error) {
	meta := map[string]interface{}{}

	annotations := MapToAttributes(metadata.Annotations)
	if len(annotations) != 0 {
		meta["annotations"] = annotations
	}

	labels := MapToAttributes(metadata.Labels)
	if len(labels) != 0 {
		meta["labels"] = labels
	}

	meta["generate_name"] = metadata.GenerateName
	meta["generation"] = metadata.Generation
	meta["name"] = metadata.Name
	meta["namespace"] = metadata.Namespace
	meta["resource_version"] = metadata.ResourceVersion
	meta["uid"] = metadata.UID
	return meta, nil
}

func CreateMetadata(metadata []interface{}) metav1.ObjectMeta {
	meta := metav1.ObjectMeta{}
	if len(metadata) < 1 {
		return meta
	}
	m := metadata[0].(map[string]interface{})

	if v, ok := m["annotations"].(map[string]interface{}); ok && len(v) > 0 {
		meta.Annotations = AttributesToMap(m["annotations"].(map[string]interface{}))
	}

	if v, ok := m["labels"].(map[string]interface{}); ok && len(v) > 0 {
		meta.Labels = AttributesToMap(m["labels"].(map[string]interface{}))
	}

	if v, ok := m["generate_name"]; ok {
		meta.GenerateName = v.(string)
	}
	if v, ok := m["name"]; ok {
		meta.Name = v.(string)
	}
	if v, ok := m["namespace"]; ok {
		meta.Namespace = v.(string)
	}

	return meta
}

func metadataFields(objectName string) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"annotations": {
			Type:        schema.TypeMap,
			Description: fmt.Sprintf("An unstructured key value map stored with the %s that may be used to store arbitrary metadata. More info: http://kubernetes.io/docs/user-guide/annotations", objectName),
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"generation": {
			Type:        schema.TypeInt,
			Description: "A sequence number representing a specific generation of the desired state.",
			Computed:    true,
		},
		"labels": {
			Type:        schema.TypeMap,
			Description: fmt.Sprintf("Map of string keys and values that can be used to organize and categorize (scope and select) the %s. May match selectors of replication controllers and services. More info: http://kubernetes.io/docs/user-guide/labels", objectName),
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"name": {
			Type:        schema.TypeString,
			Description: fmt.Sprintf("Name of the %s, must be unique. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/identifiers#names", objectName),
			Optional:    true,
			ForceNew:    true,
			Computed:    true,
		},
		"resource_version": {
			Type:        schema.TypeString,
			Description: fmt.Sprintf("An opaque value that represents the internal version of this %s that can be used by clients to determine when %s has changed. Read more: https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#concurrency-control-and-consistency", objectName, objectName),
			Computed:    true,
		},
		"uid": {
			Type:        schema.TypeString,
			Description: fmt.Sprintf("The unique in time and space value for this %s. More info: http://kubernetes.io/docs/user-guide/identifiers#uids", objectName),
			Computed:    true,
		},
	}
}
