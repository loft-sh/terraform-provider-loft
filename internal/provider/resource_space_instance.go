package provider

import (
	"context"
	storagev1 "github.com/loft-sh/api/v2/pkg/apis/storage/v1"
	client "github.com/loft-sh/loftctl/v2/pkg/client"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	agentv1 "github.com/loft-sh/agentapi/v2/pkg/apis/loft/cluster/v1"
	managementv1 "github.com/loft-sh/api/v2/pkg/apis/management/v1"
	"github.com/loft-sh/loftctl/v2/pkg/client/naming"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func resourceSpaceInstance() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "The `loft_space` resource is used to manage a Loft space.",

		CreateContext: resourceSpaceInstanceCreate,
		ReadContext:   resourceSpaceInstanceRead,
		UpdateContext: resourceSpaceInstanceUpdate,
		DeleteContext: resourceSpaceInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: spaceInstanceAttributes(),
	}
}

func resourceSpaceInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	loftClient, ok := meta.(client.Client)
	if !ok {
		return diag.Errorf("Could not access loft client")
	}

	projectName := d.Get("project").(string)

	managementClient, err := loftClient.Management()
	if err != nil {
		return diag.FromErr(err)
	}

	spaceInstance := &managementv1.SpaceInstance{
		Spec: managementv1.SpaceInstanceSpec{},
	}

	name := d.Get("name").(string)
	if name != "" {
		spaceInstance.SetName(name)
	}

	generateName := d.Get("generate_name").(string)
	if generateName != "" {
		spaceInstance.SetGenerateName(generateName)
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

	spaceInstance.SetAnnotations(annotations)

	rawLabels := d.Get("labels").(map[string]interface{})
	labels, err := attributesToMap(rawLabels)
	if err != nil {
		return diag.FromErr(err)
	}

	spaceInstance.SetLabels(labels)

	user := d.Get("owner.0.user").(string)
	team := d.Get("owner.0.team").(string)
	if user != "" && team != "" {
		return diag.Errorf("One of user or team expected.")
	}

	//templates := d.Get("templates.0")
	spaceInstance.Spec.SpaceInstanceSpec.Template = &storagev1.SpaceTemplateDefinition{}

	spaceInstance, err = managementClient.Loft().ManagementV1().SpaceInstances(naming.ProjectNamespace(projectName)).Create(ctx, spaceInstance, metav1.CreateOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	err = readSpaceInstance(projectName, spaceInstance, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSpaceInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	loftClient, ok := meta.(client.Client)
	if !ok {
		return diag.Errorf("Could not access loft client")
	}

	projectName, spaceName := parseSpaceInstanceID(d.Id())
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

	return nil
}

func resourceSpaceInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	loftClient, ok := meta.(client.Client)
	if !ok {
		return diag.Errorf("Could not access loft client")
	}

	projectName, spaceName := parseSpaceInstanceID(d.Id())
	managementClient, err := loftClient.Management()
	if err != nil {
		return diag.FromErr(err)
	}

	projectNamespace := naming.ProjectNamespace(projectName)
	oldSpace, err := managementClient.Loft().ManagementV1().SpaceInstances(projectNamespace).Get(ctx, spaceName, metav1.GetOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	modifiedSpace := oldSpace.DeepCopy()

	if d.HasChange("owner.0.user") {
		_, newUser := d.GetChange("user")
		modifiedSpace.Spec.Owner.User = newUser.(string)
	}

	if d.HasChange("owner.0.team") {
		_, newTeam := d.GetChange("team")
		modifiedSpace.Spec.Owner.Team = newTeam.(string)
	}

	//if d.HasChange("objects") {
	//	_, newObjects := d.GetChange("objects")
	//	modifiedSpace.Spec.Objects = newObjects.(string)
	//}

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

	patch := ctrlclient.MergeFrom(oldSpace)
	rawPatch, err := patch.Data(modifiedSpace)
	if err != nil {
		return diag.FromErr(err)
	}

	//space, err := clusterClient.Agent().ClusterV1().Spaces().Patch(ctx, spaceName, patch.Type(), rawPatch, metav1.PatchOptions{})
	spaceInstance, err := managementClient.Loft().ManagementV1().SpaceInstances(projectNamespace).Patch(ctx, spaceName, patch.Type(), rawPatch, metav1.PatchOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	err = readSpaceInstance(projectName, spaceInstance, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSpaceInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	loftClient, ok := meta.(client.Client)
	if !ok {
		return diag.Errorf("Could not access loft client")
	}

	projectName := d.Get("project").(string)
	spaceName := d.Get("name").(string)
	managementClient, err := loftClient.Management()
	if err != nil {
		return diag.FromErr(err)
	}

	err = managementClient.Loft().ManagementV1().SpaceInstances(naming.ProjectNamespace(projectName)).Delete(ctx, spaceName, metav1.DeleteOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
