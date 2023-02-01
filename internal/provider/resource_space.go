package provider

import (
	"context"
	client "github.com/loft-sh/loftctl/v2/pkg/client"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"

	agentv1 "github.com/loft-sh/agentapi/v2/pkg/apis/loft/cluster/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ResourceSpace() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "The `loft_space` resource is used to manage a Loft space.",

		DeprecationMessage: "`loft_space` has been deprecated and will be removed in a future release. Please use `loft_space_instance` instead.",

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
	loftClient, ok := meta.(client.Client)
	if !ok {
		return diag.Errorf("Could not access loft client")
	}

	clusterName := d.Get("cluster").(string)

	clusterClient, err := loftClient.Cluster(clusterName)
	if err != nil {
		return diag.FromErr(err)
	}

	space := &agentv1.Space{
		Spec: agentv1.SpaceSpec{},
	}

	name := d.Get("name").(string)
	if name != "" {
		space.SetName(name)
	}

	generateName := d.Get("generate_name").(string)
	if generateName != "" {
		space.SetGenerateName(generateName)
	}

	rawAnnotations := d.Get("annotations").(map[string]interface{})
	annotations, err := attributesToMap(rawAnnotations)
	if err != nil {
		return diag.FromErr(err)
	}

	sleepAfter := d.Get("sleep_after").(string)
	if sleepAfter != "" {
		duration, err := time.ParseDuration(sleepAfter)
		if err != nil {
			return diag.FromErr(err)
		}

		annotations[agentv1.SleepModeSleepAfterAnnotation] = strconv.Itoa(int(duration.Seconds()))
	}

	deleteAfter := d.Get("delete_after").(string)
	if deleteAfter != "" {
		duration, err := time.ParseDuration(deleteAfter)
		if err != nil {
			return diag.FromErr(err)
		}

		annotations[agentv1.SleepModeDeleteAfterAnnotation] = strconv.Itoa(int(duration.Seconds()))
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

	rawLabels := d.Get("labels").(map[string]interface{})
	labels, err := attributesToMap(rawLabels)
	if err != nil {
		return diag.FromErr(err)
	}

	spaceConstraints := d.Get("space_constraints").(string)
	if spaceConstraints != "" {
		labels[SpaceLabelSpaceConstraints] = spaceConstraints
	}

	space.SetLabels(labels)

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

	objects := d.Get("objects").(string)
	if objects != "" {
		space.Spec.Objects = objects
	}

	space, err = clusterClient.Agent().ClusterV1().Spaces().Create(ctx, space, metav1.CreateOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	err = readSpace(clusterName, space, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSpaceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	loftClient, ok := meta.(client.Client)
	if !ok {
		return diag.Errorf("Could not access loft client")
	}

	clusterName, spaceName := parseSpaceID(d.Id())
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

	return nil
}

func resourceSpaceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	loftClient, ok := meta.(client.Client)
	if !ok {
		return diag.Errorf("Could not access loft client")
	}

	clusterName, spaceName := parseSpaceID(d.Id())
	clusterClient, err := loftClient.Cluster(clusterName)
	if err != nil {
		return diag.FromErr(err)
	}

	oldSpace, err := clusterClient.Agent().ClusterV1().Spaces().Get(ctx, spaceName, metav1.GetOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	modifiedSpace := oldSpace.DeepCopy()

	if d.HasChange("user") {
		_, newUser := d.GetChange("user")
		modifiedSpace.Spec.User = newUser.(string)
	}

	if d.HasChange("team") {
		_, newTeam := d.GetChange("team")
		modifiedSpace.Spec.Team = newTeam.(string)
	}

	if d.HasChange("objects") {
		_, newObjects := d.GetChange("objects")
		modifiedSpace.Spec.Objects = newObjects.(string)
	}

	if d.HasChange("annotations") {
		oldAnnotations, newAnnotations := d.GetChange("annotations")

		added, modified, deleted, err := getAddedModifiedAndDeleted(
			oldAnnotations.(map[string]interface{}),
			newAnnotations.(map[string]interface{}),
		)

		if err != nil {
			return diag.FromErr(err)
		}

		for k, v := range added {
			modifiedSpace.Annotations[k] = v.(string)
		}

		for k, v := range modified {
			modifiedSpace.Annotations[k] = v.(string)
		}

		for k := range deleted {
			delete(modifiedSpace.Annotations, k)
		}
	}

	if d.HasChange("labels") {
		oldLabels, newLabels := d.GetChange("labels")

		added, modified, deleted, err := getAddedModifiedAndDeleted(
			oldLabels.(map[string]interface{}),
			newLabels.(map[string]interface{}),
		)
		if err != nil {
			return diag.FromErr(err)
		}

		for k, v := range added {
			modifiedSpace.Labels[k] = v.(string)
		}

		for k, v := range modified {
			modifiedSpace.Labels[k] = v.(string)
		}

		for k := range deleted {
			delete(modifiedSpace.Labels, k)
		}
	}

	if d.HasChange("sleep_after") {
		_, newSleepAfter := d.GetChange("sleep_after")
		sleepAfter, ok := newSleepAfter.(string)
		if !ok {
			return diag.Errorf("sleep_after value is not a string")
		}

		if sleepAfter != "" {
			duration, err := time.ParseDuration(sleepAfter)
			if err != nil {
				return diag.FromErr(err)
			}

			modifiedSpace.Annotations[agentv1.SleepModeSleepAfterAnnotation] = strconv.Itoa(int(duration.Seconds()))
		} else {
			delete(modifiedSpace.Annotations, agentv1.SleepModeSleepAfterAnnotation)
		}
	}

	if d.HasChange("delete_after") {
		_, newDeleteAfter := d.GetChange("delete_after")
		deleteAfter, ok := newDeleteAfter.(string)
		if !ok {
			return diag.Errorf("delete_after value is not an integer")
		}

		if deleteAfter != "" {
			duration, err := time.ParseDuration(deleteAfter)
			if err != nil {
				return diag.FromErr(err)
			}
			modifiedSpace.Annotations[agentv1.SleepModeDeleteAfterAnnotation] = strconv.Itoa(int(duration.Seconds()))
		} else {
			delete(modifiedSpace.Annotations, agentv1.SleepModeDeleteAfterAnnotation)
		}
	}

	if d.HasChange("sleep_schedule") {
		_, newSleepSchedule := d.GetChange("sleep_schedule")
		sleepSchedule, ok := newSleepSchedule.(string)
		if !ok {
			return diag.Errorf("sleep_schedule value is not a string")
		}

		if sleepSchedule != "" {
			modifiedSpace.Annotations[agentv1.SleepModeSleepScheduleAnnotation] = sleepSchedule
		} else {
			delete(modifiedSpace.Annotations, agentv1.SleepModeSleepScheduleAnnotation)
		}
	}

	if d.HasChange("wakeup_schedule") {
		_, newWakeupSchedule := d.GetChange("wakeup_schedule")
		wakeupSchedule, ok := newWakeupSchedule.(string)
		if !ok {
			return diag.Errorf("wakeup_schedule value is not a string")
		}

		if wakeupSchedule != "" {
			modifiedSpace.Annotations[agentv1.SleepModeWakeupScheduleAnnotation] = wakeupSchedule
		} else {
			delete(modifiedSpace.Annotations, agentv1.SleepModeWakeupScheduleAnnotation)
		}
	}

	if d.HasChange("space_constraints") {
		_, newSpaceConstraints := d.GetChange("space_constraints")
		spaceConstraints, ok := newSpaceConstraints.(string)
		if !ok {
			return diag.Errorf("space_constraints value is not a string")
		}

		if spaceConstraints != "" {
			modifiedSpace.Labels[SpaceLabelSpaceConstraints] = spaceConstraints
		} else {
			delete(modifiedSpace.Labels, SpaceLabelSpaceConstraints)
		}
	}

	patch := ctrlclient.MergeFrom(oldSpace)
	rawPatch, err := patch.Data(modifiedSpace)
	if err != nil {
		return diag.FromErr(err)
	}

	space, err := clusterClient.Agent().ClusterV1().Spaces().Patch(ctx, spaceName, patch.Type(), rawPatch, metav1.PatchOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	err = readSpace(clusterName, space, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSpaceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	loftClient, ok := meta.(client.Client)
	if !ok {
		return diag.Errorf("Could not access loft client")
	}

	clusterName := d.Get("cluster").(string)
	spaceName := d.Get("name").(string)

	clusterClient, err := loftClient.Cluster(clusterName)
	if err != nil {
		return diag.FromErr(err)
	}

	err = clusterClient.Agent().ClusterV1().Spaces().Delete(ctx, spaceName, metav1.DeleteOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
