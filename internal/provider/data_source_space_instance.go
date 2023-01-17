package provider

import (
	"context"
	"github.com/loft-sh/loftctl/v2/pkg/client"
	"github.com/loft-sh/loftctl/v2/pkg/client/naming"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSpaceInstance() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "The `loft_space_instance` data source provides information about an existing Loft space that matches the given `cluster`, `project`, and `name`.",

		ReadContext: dataSourceSpaceInstanceRead,

		Schema: spaceInstanceDataSourceAttributes(),
	}
}

func dataSourceSpaceInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	spaceName := d.Get("name").(string)
	projectName := d.Get("project").(string)

	loftClient, ok := meta.(client.Client)
	if !ok {
		return diag.Errorf("Could not access loft client")
	}

	managementClient, err := loftClient.Management()
	if err != nil {
		return diag.FromErr(err)
	}

	spaceInstance, err := managementClient.Loft().ManagementV1().SpaceInstances(naming.ProjectNamespace(projectName)).Get(ctx, spaceName, metav1.GetOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	err = readSpaceInstance(projectName, spaceInstance, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func spaceInstanceDataSourceAttributes() map[string]*schema.Schema {
	attributes := spaceInstanceAttributes()
	attributes["name"].Computed = false
	attributes["name"].Optional = false
	attributes["name"].Required = true
	attributes["name"].ConflictsWith = nil
	attributes["generate_name"].ConflictsWith = nil
	return attributes
}
