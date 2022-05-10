package provider

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	agentv1 "github.com/loft-sh/agentapi/v2/pkg/apis/loft/cluster/v1"
)

const (
	VirtualClusterLabelVirtualClusterConstraints = "loft.sh/virtual-cluster-constraints"
	DefaultVirtualClusterConstraints             = "default"
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
		"cluster": {
			// This description is used by the documentation generator and the language server.
			Description: "The cluster where the virtualCluster is located",
			Type:        schema.TypeString,
			Required:    true,
		},
		"name": {
			// This description is used by the documentation generator and the language server.
			Description: "The name of the virtualCluster",
			Type:        schema.TypeString,
			Required:    true,
		},
		"chart_name": {
			Description: "chart_name to configure chart for this virtualCluster",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"chart_version": {
			Description: "chart_version to configure chart for this virtualCluster",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"namespace": {
			// This description is used by the documentation generator and the language server.
			Description: "The namespace of the virtualCluster",
			Type:        schema.TypeString,
			Required:    true,
		},
		"annotations": {
			Description: "Annotations to configure on this virtualCluster",
			Type:        schema.TypeMap,
			Optional:    true,
		},
		"labels": {
			Description: "Labels to configure on this virtualCluster",
			Type:        schema.TypeMap,
			Optional:    true,
		},
		"virtual_cluster_constraints": {
			Description: "VirtualCluster Constraints are resources, permissions or namespace metadata that is applied and synced automatically into the virtualCluster. This is useful to ensure certain Kubernetes objects are present in each namespace to provide namespace isolation or to ensure certain labels or annotations are set on the namespace of the user.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"objects": {
			// This description is used by the documentation generator and the language server.
			Description: "Objects are Kubernetes style yamls that should get deployed into the virtualCluster",
			Type:        schema.TypeString,
			Required:    false,
			Optional:    true,
		},
	}
}

func readVirtualCluster(clusterName, namespace string, virtualCluster *agentv1.VirtualCluster, d *schema.ResourceData) error {
	virtualClusterName := virtualCluster.GetName()

	d.SetId(generateVirtualClusterId(clusterName, namespace, virtualClusterName))
	_ = d.Set("name", virtualClusterName)
	_ = d.Set("cluster", clusterName)
	_ = d.Set("namespace", namespace)
	_ = d.Set("objects", virtualCluster.Spec.Objects)
	_ = d.Set("chart_name", virtualCluster.Spec.VirtualClusterSpec.HelmRelease.Chart.Name)
	_ = d.Set("chart_version", virtualCluster.Spec.VirtualClusterSpec.HelmRelease.Chart.Version)
	safeAnnotations := removeInternalKeys(virtualCluster.GetAnnotations(), map[string]interface{}{})
	annotations, err := mapToAttributes(safeAnnotations)
	if err != nil {
		return err
	}
	_ = d.Set("annotations", annotations)

	rawLabels := virtualCluster.GetLabels()
	if rawLabels[VirtualClusterLabelVirtualClusterConstraints] != DefaultVirtualClusterConstraints {
		virtualClusterConstraints := rawLabels[VirtualClusterLabelVirtualClusterConstraints]
		_ = d.Set("virtual_cluster_constraints", virtualClusterConstraints)
	}

	safeLabels := removeInternalKeys(rawLabels, map[string]interface{}{})
	labels, err := mapToAttributes(safeLabels)
	if err != nil {
		return err
	}
	_ = d.Set("labels", labels)

	return nil
}
