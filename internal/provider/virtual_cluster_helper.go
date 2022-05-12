package provider

import (
	"fmt"
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
	err := d.Set("name", virtualClusterName)
	if err != nil {
		fmt.Println(err)
	}
	err = d.Set("cluster", clusterName)
	if err != nil {
		fmt.Println(err)
	}
	err = d.Set("namespace", namespace)
	if err != nil {
		fmt.Println(err)
	}
	err = d.Set("objects", virtualCluster.Spec.Objects)
	if err != nil {
		fmt.Println(err)
	}
	err = d.Set("chart_name", virtualCluster.Spec.VirtualClusterSpec.HelmRelease.Chart.Name)
	if err != nil {
		fmt.Println(err)
	}
	err = d.Set("chart_version", virtualCluster.Spec.VirtualClusterSpec.HelmRelease.Chart.Version)
	if err != nil {
		fmt.Println(err)
	}
	safeAnnotations := removeInternalKeys(virtualCluster.GetAnnotations(), map[string]interface{}{})
	annotations, err := mapToAttributes(safeAnnotations)
	if err != nil {
		return err
	}

	err = d.Set("annotations", annotations)
	if err != nil {
		fmt.Println(err)
	}

	rawLabels := virtualCluster.GetLabels()
	safeLabels := removeInternalKeys(rawLabels, map[string]interface{}{})
	labels, err := mapToAttributes(safeLabels)
	if err != nil {
		return err
	}
	err = d.Set("labels", labels)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}
