package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func dataSourceSpace() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "A Loft Space.",

		ReadContext: dataSourceSpaceRead,

		Schema: map[string]*schema.Schema{
			"cluster": {
				// This description is used by the documentation generator and the language server.
				Description: "The cluster where the space is located",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				// This description is used by the documentation generator and the language server.
				Description: "The name of the space",
				Type:        schema.TypeString,
				Required:    true,
			},
			"user": {
				// This description is used by the documentation generator and the language server.
				Description: "The user that owns this space",
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
			},
			"team": {
				// This description is used by the documentation generator and the language server.
				Description: "The team that owns this space",
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
			},
		},
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

	spaceId := strings.Join([]string{clusterName, space.GetName()}, "/")
	d.SetId(spaceId)

	user := space.Spec.User
	d.Set("user", user)

	team := space.Spec.Team
	d.Set("team", team)

	return diags
}
