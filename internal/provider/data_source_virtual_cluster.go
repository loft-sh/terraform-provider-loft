package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func dataSourceVirtualCluster() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "The `loft_virtual_cluster` data source provides information about an existing virtual cluster that matches the given `cluster`, `namespace`, and `name`.",

		ReadContext: dataSourceVirtualClusterRead,

		Schema: virtualClustersAttributes(),
	}
}

func dataSourceVirtualClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	virtualClusterName := d.Get("name").(string)
	clusterName := d.Get("cluster").(string)
	namespace := d.Get("namespace").(string)
	apiClient, ok := meta.(*apiClient)
	if !ok {
		return diag.Errorf("Could not access apiClient")
	}

	clusterClient, err := apiClient.LoftClient.Cluster(clusterName)
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
	schema := virtualClusterAttributes()
	schema["name"].Computed = false
	schema["name"].ConflictsWith = nil
	schema["name"].Optional = false
	schema["name"].Required = true
	schema["generate_name"].ConflictsWith = nil
	return schema
}
