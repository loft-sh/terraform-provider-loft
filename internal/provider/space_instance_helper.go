package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	v1 "github.com/loft-sh/api/v2/pkg/apis/management/v1"
	"strings"
)

func spaceInstanceAttributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Unique identifier for this space. The format is `<project>/<name>`.",
		},
		"project": {
			// This description is used by the documentation generator and the language server.
			Description: "The project to which the space belongs.",
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
		"owner": {
			Description: "The owner of the space.",
			Type:        schema.TypeList,
			MinItems:    1,
			MaxItems:    1,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"team": {
						Description: "The team that owns this space.",
						Type:        schema.TypeString,
						Required:    false,
						Optional:    true,
					},
					"user": {
						Description: "The user that owns this space.",
						Type:        schema.TypeString,
						Required:    false,
						Optional:    true,
					},
				},
			},
		},
		"template": {
			// This description is used by the documentation generator and the language server.
			Description: "The inline template to use for space creation. This is mutually exclusive with templateRef.",
			Type:        schema.TypeList,
			MinItems:    1,
			MaxItems:    1,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					/*"metadata": &schema.Schema{
						MinItems: 1,
						MaxItems: 1,
						Optional: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
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
							},
						},
					},*/
					"objects": {
						// This description is used by the documentation generator and the language server.
						Description: "Objects are Kubernetes style yamls that should get deployed into the space.",
						Type:        schema.TypeString,
						Required:    false,
						Optional:    true,
					},
					"charts": {
						Type:     schema.TypeList,
						Optional: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Description: "Name is the chart name in the repository",
									Type:        schema.TypeString,
									Required:    false,
									Optional:    true,
								},
								"version": {
									Description: "Version is the chart version in the repository",
									Type:        schema.TypeString,
									Required:    false,
									Optional:    true,
								},
								"repo_url": {
									Description: "RepoURL is the repo url where the chart can be found",
									Type:        schema.TypeString,
									Required:    false,
									Optional:    true,
								},
								"username": {
									Description: "The username that is required for this repository",
									Type:        schema.TypeString,
									Required:    false,
									Optional:    true,
								},
								"password": {
									Description: "The password that is required for this repository",
									Type:        schema.TypeString,
									Required:    false,
									Optional:    true,
								},
								"insecure_skip_tls_verify": {
									Description: "If tls certificate checks for the chart download should be skipped",
									Type:        schema.TypeBool,
									Required:    false,
									Optional:    true,
								},
								"release_name": {
									Description: "ReleaseName is the preferred release name of the app",
									Type:        schema.TypeString,
									Required:    false,
									Optional:    true,
								},
								"release_namespace": {
									Description: "ReleaseNamespace is the preferred release namespace of the app",
									Type:        schema.TypeString,
									Required:    false,
									Optional:    true,
								},
								"values": {
									Description: "Values are the values that should get passed to the chart",
									Type:        schema.TypeString,
									Required:    false,
									Optional:    true,
								},
								"wait": {
									Description: "Wait determines if Loft should wait during deploy for the app to become ready",
									Type:        schema.TypeBool,
									Required:    false,
									Optional:    true,
								},
								"timeout": {
									Description:      "Timeout is the time to wait for any individual Kubernetes operation (like Jobs for hooks) (default 5m0s)",
									Type:             schema.TypeString,
									Optional:         true,
									StateFunc:        durationToSeconds,
									ValidateDiagFunc: validateDuration,
								},
							},
						},
					},
					"apps": {
						Type:     schema.TypeList,
						Optional: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Description: "Name of the target app",
									Type:        schema.TypeString,
									Required:    false,
									Optional:    true,
								},
								"namespace": {
									Description: "Namespace specifies in which target namespace the app should get deployed in",
									Type:        schema.TypeString,
									Required:    false,
									Optional:    true,
								},
								"release_name": {
									Description: "ReleaseName is the name of the app release",
									Type:        schema.TypeString,
									Required:    false,
									Optional:    true,
								},
								"version": {
									Description: "Version of the app",
									Type:        schema.TypeString,
									Required:    false,
									Optional:    true,
								},
								"parameters": {
									Description: "Parameters to use for the app",
									Type:        schema.TypeString,
									Required:    false,
									Optional:    true,
								},
							},
						},
					},
					"access": {
						Type:     schema.TypeList,
						MinItems: 1,
						MaxItems: 1,
						Optional: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"default_cluster_role": {
									Description: "Specifies which cluster role should get applied to users or teams that do not match a rule below.",
									Type:        schema.TypeString,
									Optional:    true,
								},
								"rules": {
									Description: "Rules defines which users and teams should have which access to the space. If no rule matches an authenticated incoming user, the user will get cluster admin access.",
									Type:        schema.TypeList,
									Optional:    true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"users": {
												Description: "Users this rule matches. * means all users.",
												Type:        schema.TypeList,
												Optional:    true,
												Elem: &schema.Schema{
													Type: schema.TypeString,
												},
											},
											"teams": {
												Description: "Teams that this rule matches. * means all teams.",
												Type:        schema.TypeList,
												Optional:    true,
												Elem: &schema.Schema{
													Type: schema.TypeString,
												},
											},
											"cluster_role": {
												Description: "ClusterRole is the cluster role that should be assigned to the member.",
												Type:        schema.TypeString,
												Required:    false,
												Optional:    true,
											},
										},
									},
								},
							},
						},
					},
				},
			},
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
	}
}

func readSpaceInstance(projectName string, space *v1.SpaceInstance, d *schema.ResourceData) error {
	return nil
}

func parseSpaceInstanceID(id string) (projectName, spaceName string) {
	projectName = ""
	spaceName = ""

	tokens := strings.Split(id, "/")
	if len(tokens) == 2 {
		projectName = tokens[0]
		spaceName = tokens[1]
	}

	return
}
