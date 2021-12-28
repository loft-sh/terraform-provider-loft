package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	agentv1 "github.com/loft-sh/agentapi/v2/pkg/apis/loft/cluster/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func resourceSpace() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "A Loft Space.",

		CreateContext: resourceSpaceCreate,
		ReadContext:   resourceSpaceRead,
		UpdateContext: resourceSpaceUpdate,
		DeleteContext: resourceSpaceDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				// This description is used by the documentation generator and the language server.
				Description: "The name of the space",
				Type:        schema.TypeString,
				Required:    true,
			},
			"cluster": {
				// This description is used by the documentation generator and the language server.
				Description: "The cluster where the space is managed",
				Type:        schema.TypeString,
				Required:    true,
			},
			"annotations": {
				Description: "Annotations to configure on this space",
				Type:        schema.TypeMap,
				Optional:    true,
			},
			"labels": {
				Description: "Labels to configure on this space",
				Type:        schema.TypeMap,
				Optional:    true,
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

func resourceSpaceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)
	apiClient, ok := meta.(*apiClient)
	if !ok {
		return diag.Errorf("Could not access apiClient")
	}

	clusterName := d.Get("cluster").(string)
	spaceName := d.Get("name").(string)

	clusterClient, err := apiClient.LoftClient.Cluster(clusterName)
	if err != nil {
		return diag.FromErr(err)
	}

	user := d.Get("user").(string)
	team := d.Get("team").(string)
	if user != "" && team != "" {
		return diag.Errorf("One of user or team expected.")
	}

	space := &agentv1.Space{
		Spec: agentv1.SpaceSpec{},
	}
	space.SetName(spaceName)

	annotations := d.Get("annotations").(map[string]interface{})
	if len(annotations) > 0 {
		strAnnotations, err := flattenMap(annotations)
		if err != nil {
			return diag.FromErr(err)
		}
		space.SetAnnotations(strAnnotations)
	}

	labels := d.Get("labels").(map[string]interface{})
	if len(labels) > 0 {
		strLabels, err := flattenMap(labels)
		if err != nil {
			return diag.FromErr(err)
		}
		space.SetLabels(strLabels)
	}

	if user != "" {
		space.Spec.User = user
	}

	if team != "" {
		space.Spec.Team = team
	}

	space, err = clusterClient.Agent().ClusterV1().Spaces().Create(ctx, space, metav1.CreateOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(generateSpaceId(clusterName, space.GetName()))
	d.Set("user", space.Spec.User)
	d.Set("team", space.Spec.Team)

	return nil
}

func resourceSpaceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	d.SetId(generateSpaceId(clusterName, space.GetName()))
	d.Set("user", space.Spec.User)
	d.Set("team", space.Spec.Team)

	return diags
}

func resourceSpaceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	return diag.Errorf("update not implemented")
}

func resourceSpaceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient, ok := meta.(*apiClient)
	if !ok {
		return diag.Errorf("Could not access apiClient")
	}

	clusterName := d.Get("cluster").(string)
	spaceName := d.Get("name").(string)

	clusterClient, err := apiClient.LoftClient.Cluster(clusterName)
	if err != nil {
		return diag.FromErr(err)
	}

	err = clusterClient.Agent().ClusterV1().Spaces().Delete(context.TODO(), spaceName, metav1.DeleteOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func generateSpaceId(clusterName, spaceName string) string {
	return strings.Join([]string{clusterName, spaceName}, "/")
}

func flattenMap(rawMap map[string]interface{}) (map[string]string, error) {
	strMap := map[string]string{}
	for k, v := range rawMap {
		str, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("non-string value used in map")
		}
		strMap[k] = str
	}
	return strMap, nil
}
