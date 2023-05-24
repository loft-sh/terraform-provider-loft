package provider

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	agentv1 "github.com/loft-sh/agentapi/v3/pkg/apis/loft/cluster/v1"
)

const (
	SpaceLabelSpaceConstraints = "loft.sh/space-constraints"
	DefaultSpaceConstraints    = "default"
)

func generateSpaceID(clusterName, spaceName string) string {
	return strings.Join([]string{clusterName, spaceName}, "/")
}

func parseSpaceID(id string) (clusterName, spaceName string) {
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
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Unique identifier for this space. The format is `<cluster>/<name>`.",
		},
		"cluster": {
			// This description is used by the documentation generator and the language server.
			Description: "The cluster where the space is located.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"name": {
			// This description is used by the documentation generator and the language server.
			Description:   "The name of the space.",
			Type:          schema.TypeString,
			Optional:      true,
			ForceNew:      true,
			Computed:      true,
			ConflictsWith: []string{"generate_name"},
		},
		"generate_name": {
			// This description is used by the documentation generator and the language server.
			Description:   "Prefix, used by the server, to generate a unique name ONLY IF the `name` field has not been provided. This value will also be combined with a unique suffix. Read more about [kubernetes API conventions for idempotency here](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#idempotency).",
			Type:          schema.TypeString,
			Optional:      true,
			ForceNew:      true,
			Computed:      true,
			ConflictsWith: []string{"name"},
		},
		"annotations": {
			Description: "The annotations to configure on this space.",
			Type:        schema.TypeMap,
			Optional:    true,
		},
		"labels": {
			Description: "The labels to configure on this space.",
			Type:        schema.TypeMap,
			Optional:    true,
		},
		"sleep_after": {
			// This description is used by the documentation generator and the language server.
			Description:      "If configured, this will tell Loft to put the space to sleep after the specified duration of inactivity. The format is a string accepted by the [time.ParseDuration](https://pkg.go.dev/time#ParseDuration) function, such as `\"1h\"`",
			Type:             schema.TypeString,
			Optional:         true,
			StateFunc:        durationToSeconds,
			ValidateDiagFunc: validateDuration,
		},
		"delete_after": {
			// This description is used by the documentation generator and the language server.
			Description:      "If configured, this will tell Loft to delete the space after the specified duration of inactivity. The format is a string accepted by the [time.ParseDuration](https://pkg.go.dev/time#ParseDuration) function, such as `\"1h\"`",
			Type:             schema.TypeString,
			Optional:         true,
			StateFunc:        durationToSeconds,
			ValidateDiagFunc: validateDuration,
		},
		"sleep_schedule": {
			Description: "Put the space to sleep at certain times. See [crontab.guru](https://crontab.guru/) for valid configurations. This might be useful if you want to set the space sleeping over the weekend for example.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"wakeup_schedule": {
			Description: "Wake up the space at certain times. See [crontab.guru](https://crontab.guru/) for valid configurations. This might be useful if it started sleeping due to inactivity and you want to wake up the space on a regular basis.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"space_constraints": {
			Description: "Space Constraints are resources, permissions or namespace metadata that is applied and synced automatically into the space. This is useful to ensure certain Kubernetes objects are present in each namespace to provide namespace isolation or to ensure certain labels or annotations are set on the namespace of the user.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"user": {
			Description: "The user that owns this space.",
			Type:        schema.TypeString,
			Required:    false,
			Optional:    true,
		},
		"team": {
			Description: "The team that owns this space.",
			Type:        schema.TypeString,
			Required:    false,
			Optional:    true,
		},
		"objects": {
			// This description is used by the documentation generator and the language server.
			Description: "Objects are Kubernetes style yamls that should get deployed into the space.",
			Type:        schema.TypeString,
			Required:    false,
			Optional:    true,
		},
	}
}

func readSpace(clusterName string, space *agentv1.Space, d *schema.ResourceData) error {
	spaceName := space.GetName()

	d.SetId(generateSpaceID(clusterName, spaceName))
	if err := d.Set("cluster", clusterName); err != nil {
		return err
	}
	if err := d.Set("name", spaceName); err != nil {
		return err
	}
	if err := d.Set("generate_name", space.GetGenerateName()); err != nil {
		return err
	}
	if err := d.Set("user", space.Spec.User); err != nil {
		return err
	}
	if err := d.Set("team", space.Spec.Team); err != nil {
		return err
	}
	if err := d.Set("objects", space.Spec.Objects); err != nil {
		return err
	}

	rawAnnotations := space.GetAnnotations()
	if rawAnnotations[agentv1.SleepModeSleepAfterAnnotation] != "" {
		sleepAfter := rawAnnotations[agentv1.SleepModeSleepAfterAnnotation]
		if err := d.Set("sleep_after", sleepAfter); err != nil {
			return err
		}
	}

	if rawAnnotations[agentv1.SleepModeDeleteAfterAnnotation] != "" {
		deleteAfter := rawAnnotations[agentv1.SleepModeDeleteAfterAnnotation]
		if err := d.Set("delete_after", deleteAfter); err != nil {
			return err
		}
	}

	if rawAnnotations[agentv1.SleepModeSleepScheduleAnnotation] != "" {
		sleepSchedule := rawAnnotations[agentv1.SleepModeSleepScheduleAnnotation]
		if err := d.Set("sleep_schedule", sleepSchedule); err != nil {
			return err
		}
	}

	if rawAnnotations[agentv1.SleepModeWakeupScheduleAnnotation] != "" {
		wakeupSchedule := rawAnnotations[agentv1.SleepModeWakeupScheduleAnnotation]
		if err := d.Set("wakeup_schedule", wakeupSchedule); err != nil {
			return err
		}
	}

	safeAnnotations := removeInternalKeys(space.GetAnnotations(), map[string]interface{}{})
	annotations, err := mapToAttributes(safeAnnotations)
	if err != nil {
		return err
	}
	if err := d.Set("annotations", annotations); err != nil {
		return err
	}

	rawLabels := space.GetLabels()
	if rawLabels[SpaceLabelSpaceConstraints] != DefaultSpaceConstraints {
		spaceConstraints := rawLabels[SpaceLabelSpaceConstraints]
		if err := d.Set("space_constraints", spaceConstraints); err != nil {
			return err
		}
	}

	safeLabels := removeInternalKeys(rawLabels, map[string]interface{}{})
	labels, err := mapToAttributes(safeLabels)
	if err != nil {
		return err
	}
	if err := d.Set("labels", labels); err != nil {
		return err
	}

	return nil
}

func validateDuration(v interface{}, _ cty.Path) diag.Diagnostics {
	valStr := v.(string)

	_, err := time.ParseDuration(valStr)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func durationToSeconds(val interface{}) string {
	valStr, ok := val.(string)
	if !ok {
		return ""
	}

	duration, err := time.ParseDuration(valStr)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%d", int(duration.Seconds()))
}
