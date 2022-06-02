package provider

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	agentv1 "github.com/loft-sh/agentapi/v2/pkg/apis/loft/cluster/v1"
)

func generateVirtualClusterId(clusterName, namespace, virtualClusterName string) string {
	return strings.Join([]string{clusterName, namespace, virtualClusterName}, "/")
}

func parseVirtualClusterId(id string) (clusterName, namespace, virtualClusterName string) {
	clusterName = ""
	namespace = ""
	virtualClusterName = ""

	tokens := strings.Split(id, "/")
	if len(tokens) == 3 {
		clusterName = tokens[0]
		namespace = tokens[1]
		virtualClusterName = tokens[2]
	}

	return
}

func virtualClusterAttributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Unique identifier for this virtual cluster. The format is `<cluster>/<namespace>/<name>`.",
		},
		"cluster": {
			Description: "The cluster where the virtual cluster is located.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"name": {
			Description:   "The name of the virtual cluster.",
			Type:          schema.TypeString,
			Optional:      true,
			ForceNew:      true,
			Computed:      true,
			ConflictsWith: []string{"generate_name"},
		},
		"generate_name": {
			Description:   "Prefix, used by the server, to generate a unique name ONLY IF the `name` field has not been provided. This value will also be combined with a unique suffix. Read more: https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#idempotency.",
			Type:          schema.TypeString,
			Optional:      true,
			ForceNew:      true,
			Computed:      true,
			ConflictsWith: []string{"name"},
		},
		"chart_name": {
			Description: "The helm chart name used to configure the virtual cluster. ",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"chart_version": {
			Description: "The helm chart version used to configure the virtual cluster.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"values": {
			Description: "The helm chart values to configure the virtual cluster.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"namespace": {
			Description: "The namespace where the virtual cluster is deployed.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"annotations": {
			Description: "The annotations to configure on this virtual cluster.",
			Type:        schema.TypeMap,
			Optional:    true,
		},
		"labels": {
			Description: "The labels to configure on this virtual cluster.",
			Type:        schema.TypeMap,
			Optional:    true,
		},
		"objects": {
			// This description is used by the documentation generator and the language server.
			Description: "Objects are Kubernetes style yamls that should get deployed into the virtual cluster",
			Type:        schema.TypeString,
			Required:    false,
			Optional:    true,
		},
	}
}

func readVirtualCluster(clusterName, namespace string, virtualCluster *agentv1.VirtualCluster, d *schema.ResourceData) error {
	virtualClusterName := virtualCluster.GetName()

	d.SetId(generateVirtualClusterId(clusterName, namespace, virtualClusterName))

	if err := d.Set("name", virtualClusterName); err != nil {
		return err
	}
	if err := d.Set("generate_name", virtualCluster.GetGenerateName()); err != nil {
		return err
	}
	if err := d.Set("cluster", clusterName); err != nil {
		return err
	}
	if err := d.Set("namespace", namespace); err != nil {
		return err
	}
	if err := d.Set("objects", virtualCluster.Spec.Objects); err != nil {
		return err
	}
	if err := d.Set("chart_name", virtualCluster.Spec.VirtualClusterSpec.HelmRelease.Chart.Name); err != nil {
		return err
	}
	if err := d.Set("chart_version", virtualCluster.Spec.VirtualClusterSpec.HelmRelease.Chart.Version); err != nil {
		return err
	}
	if err := d.Set("values", virtualCluster.Spec.VirtualClusterSpec.HelmRelease.Values); err != nil {
		return err
	}

	safeAnnotations := removeInternalKeys(virtualCluster.GetAnnotations(), map[string]interface{}{})
	annotations, err := mapToAttributes(safeAnnotations)
	if err != nil {
		return err
	}
	if err := d.Set("annotations", annotations); err != nil {
		return err
	}

	rawLabels := virtualCluster.GetLabels()
	safeLabels := removeInternalKeys(rawLabels, map[string]interface{}{})
	labels, err := mapToAttributes(safeLabels)
	if err != nil {
		return err
	}

	if err := d.Set("labels", labels); err != nil {
		return err
	}
	return nil
}
