package provider

import (
	"context"

	"github.com/loft-sh/loftctl/v3/pkg/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DataSourceSpace() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "The `loft_space` data source provides information about an existing Loft space that matches the given `cluster` and `name`.",

		DeprecationMessage: "`loft_space` has been deprecated and will be removed in a future release. Please use `loft_space_instance` instead.",

		ReadContext: dataSourceSpaceRead,

		Schema: spaceDataSourceAttributes(),
	}
}

func dataSourceSpaceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	spaceName := d.Get("name").(string)
	clusterName := d.Get("cluster").(string)

	loftClient, ok := meta.(client.Client)
	if !ok {
		return diag.Errorf("Could not access loft client")
	}

	clusterClient, err := loftClient.Cluster(clusterName)
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
	attributes := spaceAttributes()
	attributes["name"].Computed = false
	attributes["name"].Optional = false
	attributes["name"].Required = true
	attributes["name"].ConflictsWith = nil
	attributes["generate_name"].ConflictsWith = nil
	return attributes
}
