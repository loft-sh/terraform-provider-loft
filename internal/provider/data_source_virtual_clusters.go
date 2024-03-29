package provider

import (
	"context"
	"strings"

	"github.com/loft-sh/loftctl/v3/pkg/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	v1 "github.com/loft-sh/agentapi/v3/pkg/apis/loft/storage/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DataSourceVirtualClusters() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "The `loft_virtual_clusters` data source provides information about all virtual clusters that match the given `cluster` and `namespace`.",

		DeprecationMessage: "`loft_virtual_clusters` has been deprecated and will be removed in a future release.",

		ReadContext: dataSourceVirtualClustersRead,

		Schema: map[string]*schema.Schema{
			"virtual_clusters": {
				Description: "All virtual_clusters",
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        DataSourceVirtualCluster(),
			},
			"cluster": {
				Description: "The cluster to list virtual_clusters from.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"namespace": {
				Description: "The namespace to list virtual_clusters from.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceVirtualClustersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	clusterName := d.Get("cluster").(string)
	namespace := d.Get("namespace").(string)

	loftClient, ok := meta.(client.Client)
	if !ok {
		return diag.Errorf("Could not access loft client")
	}

	clusterClient, err := loftClient.Cluster(clusterName)
	if err != nil {
		return diag.FromErr(err)
	}

	virtualClustersList, err := clusterClient.Agent().StorageV1().VirtualClusters(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	var virtualClusters []map[string]interface{}
	for _, virtualCluster := range virtualClustersList.Items {
		flattenedVirtualCluster, err := flattenVirtualCluster(clusterName, namespace, virtualCluster)
		if err != nil {
			return diag.FromErr(err)
		}
		virtualClusters = append(virtualClusters, flattenedVirtualCluster)
	}

	virtualClusterId := strings.Join([]string{clusterName, namespace, "virtual_clusters"}, "/")
	d.SetId(virtualClusterId)

	if err := d.Set("virtual_clusters", virtualClusters); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func flattenVirtualCluster(clusterName, namespace string, virtualCluster v1.VirtualCluster) (map[string]interface{}, error) {
	flattenedVirtualCluster := map[string]interface{}{
		"name":      virtualCluster.GetName(),
		"cluster":   clusterName,
		"namespace": namespace,
		"objects":   virtualCluster.Spec.Objects,
	}

	rawAnnotations := removeInternalKeys(virtualCluster.GetAnnotations(), map[string]interface{}{})
	annotations, err := mapToAttributes(rawAnnotations)
	if err != nil {
		return nil, err
	}

	flattenedVirtualCluster["annotations"] = annotations

	return flattenedVirtualCluster, nil
}
