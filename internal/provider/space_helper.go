package provider

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	agentv1 "github.com/loft-sh/agentapi/v2/pkg/apis/loft/cluster/v1"
)

const (
	SpaceLabelSpaceConstraints = "loft.sh/space-constraints"
	DefaultSpaceConstraints    = "default"
)

func generateSpaceId(clusterName, spaceName string) string {
	return strings.Join([]string{clusterName, spaceName}, "/")
}

func parseSpaceId(id string) (clusterName, spaceName string) {
	clusterName = ""
	spaceName = ""

	tokens := strings.Split(id, "/")
	if len(tokens) == 2 {
		clusterName = tokens[0]
		spaceName = tokens[1]
	}

	return
}

func spaceAttributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
		"sleep_after": {
			// This description is used by the documentation generator and the language server.
			Description: "If set to non zero, will tell the space to sleep after specified seconds of inactivity",
			Type:        schema.TypeInt,
			Optional:    true,
		},
		"delete_after": {
			// This description is used by the documentation generator and the language server.
			Description: "If set to non zero, will tell loft to delete the space after specified seconds of inactivity",
			Type:        schema.TypeInt,
			Optional:    true,
		},
		"sleep_schedule": {
			Description: "Put the space to sleep at certain times. See crontab.guru for valid configurations. This might be useful if you want to set the space sleeping over the weekend for example.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"wakeup_schedule": {
			Description: "Wake up the space at certain times. See crontab.guru for valid configurations. This might be useful if it started sleeping due to inactivity and you want to wake up the space on a regular basis.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"space_constraints": {
			Description: "Space Constraints are resources, permissions or namespace metadata that is applied and synced automatically into the space. This is useful to ensure certain Kubernetes objects are present in each namespace to provide namespace isolation or to ensure certain labels or annotations are set on the namespace of the user.",
			Type:        schema.TypeString,
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
		"objects": {
			// This description is used by the documentation generator and the language server.
			Description: "Objects are Kubernetes style yamls that should get deployed into the space",
			Type:        schema.TypeString,
			Required:    false,
			Optional:    true,
		},
	}
}

func readSpace(clusterName string, space *agentv1.Space, d *schema.ResourceData) error {
	spaceName := space.GetName()

	d.SetId(generateSpaceId(clusterName, spaceName))
	err := d.Set("name", spaceName)
	if err != nil {
		fmt.Println(err)
	}
	err = d.Set("cluster", clusterName)
	if err != nil {
		fmt.Println(err)
	}
	err = d.Set("name", spaceName)
	if err != nil {
		fmt.Println(err)
	}
	err = d.Set("user", space.Spec.User)
	if err != nil {
		fmt.Println(err)
	}
	err = d.Set("team", space.Spec.Team)
	if err != nil {
		fmt.Println(err)
	}
	err = d.Set("objects", space.Spec.Objects)
	if err != nil {
		fmt.Println(err)
	}

	rawAnnotations := space.GetAnnotations()
	if rawAnnotations[agentv1.SleepModeSleepAfterAnnotation] != "" {
		sleepAfter, err := strconv.Atoi(rawAnnotations[agentv1.SleepModeSleepAfterAnnotation])
		if err != nil {
			return err
		}
		err = d.Set("sleep_after", sleepAfter)
		if err != nil {
			fmt.Println(err)
		}
	}

	if rawAnnotations[agentv1.SleepModeDeleteAfterAnnotation] != "" {
		deleteAfter, err := strconv.Atoi(rawAnnotations[agentv1.SleepModeDeleteAfterAnnotation])
		if err != nil {
			return err
		}
		err = d.Set("delete_after", deleteAfter)
		if err != nil {
			fmt.Println(err)
		}
	}

	if rawAnnotations[agentv1.SleepModeSleepScheduleAnnotation] != "" {
		sleepSchedule := rawAnnotations[agentv1.SleepModeSleepScheduleAnnotation]
		err = d.Set("sleep_schedule", sleepSchedule)
		if err != nil {
			fmt.Println(err)
		}
	}

	if rawAnnotations[agentv1.SleepModeWakeupScheduleAnnotation] != "" {
		wakeupSchedule := rawAnnotations[agentv1.SleepModeWakeupScheduleAnnotation]
		err = d.Set("wakeup_schedule", wakeupSchedule)
		if err != nil {
			fmt.Println(err)
		}
	}

	safeAnnotations := removeInternalKeys(space.GetAnnotations(), map[string]interface{}{})
	annotations, err := mapToAttributes(safeAnnotations)
	if err != nil {
		return err
	}
	err = d.Set("annotations", annotations)
	if err != nil {
		fmt.Println(err)
	}

	rawLabels := space.GetLabels()
	if rawLabels[SpaceLabelSpaceConstraints] != DefaultSpaceConstraints {
		spaceConstraints := rawLabels[SpaceLabelSpaceConstraints]
		err = d.Set("space_constraints", spaceConstraints)
		if err != nil {
			fmt.Println(err)
		}
	}

	safeLabels := removeInternalKeys(rawLabels, map[string]interface{}{})
	labels, err := mapToAttributes(safeLabels)
	if err != nil {
		return err
	}
	err = d.Set("labels", labels)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}
