package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func dataSourceSpace() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "The `loft_space` data source provides information about an existing Loft space that matches the given `cluster` and `name`.",

		ReadContext: dataSourceSpaceRead,

		Schema: spaceDataSourceAttributes(),
	}
}

func dataSourceSpaceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	spaceName := d.Get("name").(string)
	clusterName := d.Get("cluster").(string)

	apiClient, ok := meta.(*apiClient)
	if !ok {
		return diag.Errorf("Could not access apiClient")
	}

	clusterClient, err := apiClient.LoftClient.Cluster(clusterName)
	if err != nil {
		return diag.FromErr(err)
	}

	space, err := clusterClient.Agent().ClusterV1().Spaces().Get(ctx, spaceName, metav1.GetOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	err = readSpace(clusterName, space, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func spaceDataSourceAttributes() map[string]*schema.Schema {
	schema := spaceAttributes()
	schema["name"].Computed = false
	schema["name"].Optional = false
	schema["name"].Required = true
	return schema
}
