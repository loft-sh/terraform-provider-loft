package provider

import (
	"context"
	"strconv"

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
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: spaceAttributes(),
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

	space := &agentv1.Space{
		Spec: agentv1.SpaceSpec{},
	}
	space.SetName(spaceName)

	rawAnnotations := d.Get("annotations").(map[string]interface{})
	annotations := map[string]string{}
	if len(rawAnnotations) > 0 {
		annotations, err = attributesToMap(rawAnnotations)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	sleepAfter := d.Get("sleep_after").(int)
	if sleepAfter > 0 {
		annotations[agentv1.SleepModeSleepAfterAnnotation] = strconv.Itoa(sleepAfter)
	}

	deleteAfter := d.Get("delete_after").(int)
	if deleteAfter > 0 {
		annotations[agentv1.SleepModeDeleteAfterAnnotation] = strconv.Itoa(deleteAfter)
	}

	sleepSchedule := d.Get("sleep_schedule").(string)
	if sleepSchedule != "" {
		annotations[agentv1.SleepModeSleepScheduleAnnotation] = sleepSchedule
	}

	wakeupSchedule := d.Get("wakeup_schedule").(string)
	if wakeupSchedule != "" {
		annotations[agentv1.SleepModeWakeupScheduleAnnotation] = wakeupSchedule
	}

	space.SetAnnotations(annotations)

	labels := d.Get("labels").(map[string]interface{})
	if len(labels) > 0 {
		strLabels, err := attributesToMap(labels)
		if err != nil {
			return diag.FromErr(err)
		}
		space.SetLabels(strLabels)
	}

	user := d.Get("user").(string)
	team := d.Get("team").(string)
	if user != "" && team != "" {
		return diag.Errorf("One of user or team expected.")
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

	apiClient, ok := meta.(*apiClient)
	if !ok {
		return diag.Errorf("Could not access apiClient")
	}

	clusterName, spaceName := parseSpaceId(d.Id())
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
