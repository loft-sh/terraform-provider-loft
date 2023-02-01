package provider

import (
	"context"
	"github.com/loft-sh/loftctl/v2/pkg/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DataSourceVirtualCluster() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "The `loft_virtual_cluster` data source provides information about an existing virtual cluster that matches the given `cluster`, `namespace`, and `name`.",

		DeprecationMessage: "`loft_virtual_cluster` has been deprecated and will be removed in a future release. Please use `loft_virtual_cluster_instance` instead.",

		ReadContext: dataSourceVirtualClusterRead,

		Schema: virtualClustersAttributes(),
	}
}

func dataSourceVirtualClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	virtualClusterName := d.Get("name").(string)
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

	virtualCluster, err := clusterClient.Agent().StorageV1().VirtualClusters(namespace).Get(ctx, virtualClusterName, metav1.GetOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	if err := readVirtualCluster(clusterName, namespace, virtualCluster, d); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func virtualClustersAttributes() map[string]*schema.Schema {
	attributes := virtualClusterAttributes()
	attributes["name"].Computed = false
	attributes["name"].ConflictsWith = nil
	attributes["name"].Optional = false
	attributes["name"].Required = true
	attributes["generate_name"].ConflictsWith = nil
	return attributes
}
