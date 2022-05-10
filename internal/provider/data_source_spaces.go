package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	v1 "github.com/loft-sh/agentapi/v2/pkg/apis/loft/cluster/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func dataSourceSpaces() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "A Loft Space.",

		ReadContext: dataSourceSpacesRead,

		Schema: map[string]*schema.Schema{
			"spaces": {
				Description: "All spaces, or the spaces matching the given label selector",
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataSourceSpace(),
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

	apiClient, ok := meta.(*apiClient)
	if !ok {
		return diag.Errorf("Could not access apiClient")
	}

	clusterClient, err := apiClient.LoftClient.Cluster(clusterName)
	if err != nil {
		return diag.FromErr(err)
	}

	spacesList, err := clusterClient.Agent().ClusterV1().Spaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	spaces := []map[string]interface{}{}
	for _, space := range spacesList.Items {
		flattenedSpace, err := flattenSpace(clusterName, space)
		if err != nil {
			return diag.FromErr(err)
		}
		spaces = append(spaces, flattenedSpace)
	}

	spaceId := strings.Join([]string{clusterName, "spaces"}, "/")
	d.SetId(spaceId)
	_ = d.Set("spaces", spaces)
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
