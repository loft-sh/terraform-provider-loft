package provider

import (
	"context"
	"strings"

	"github.com/loft-sh/loftctl/v3/pkg/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	v1 "github.com/loft-sh/agentapi/v3/pkg/apis/loft/cluster/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DataSourceSpaces() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "The `loft_spaces` data source provides information about all Loft spaces in the given `cluster`.",

		DeprecationMessage: "`loft_spaces` has been deprecated and will be removed in a future release.",

		ReadContext: dataSourceSpacesRead,

		Schema: map[string]*schema.Schema{
			"spaces": {
				Description: "All spaces",
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        DataSourceSpace(),
			},
			"cluster": {
				Description: "The cluster to list spaces from.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceSpacesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	clusterName := d.Get("cluster").(string)

	loftClient, ok := meta.(client.Client)
	if !ok {
		return diag.Errorf("Could not access loft client")
	}

	clusterClient, err := loftClient.Cluster(clusterName)
	if err != nil {
		return diag.FromErr(err)
	}

	spacesList, err := clusterClient.Agent().ClusterV1().Spaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	var spaces []map[string]interface{}
	for _, space := range spacesList.Items {
		flattenedSpace, err := flattenSpace(clusterName, space)
		if err != nil {
			return diag.FromErr(err)
		}
		spaces = append(spaces, flattenedSpace)
	}

	spaceID := strings.Join([]string{clusterName, "spaces"}, "/")
	d.SetId(spaceID)
	if err := d.Set("spaces", spaces); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func flattenSpace(clusterName string, space v1.Space) (map[string]interface{}, error) {
	flattenedSpace := map[string]interface{}{
		"name":    space.GetName(),
		"cluster": clusterName,
		"user":    space.Spec.User,
		"team":    space.Spec.Team,
		"objects": space.Spec.Objects,
	}

	rawAnnotations := removeInternalKeys(space.GetAnnotations(), map[string]interface{}{})
	annotations, err := mapToAttributes(rawAnnotations)
	if err != nil {
		return nil, err
	}

	flattenedSpace["annotations"] = annotations

	return flattenedSpace, nil
}
