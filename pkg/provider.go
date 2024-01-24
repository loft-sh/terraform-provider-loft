package loft

import (
	"context"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/loft-sh/loftctl/v3/pkg/client"
	legacy "github.com/loft-sh/terraform-provider-loft/internal/provider"
	"github.com/loft-sh/terraform-provider-loft/pkg/resources"
)

func init() {
	// Set descriptions to support Markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown
}

func New() func() *schema.Provider {
	return func() *schema.Provider {
		return &schema.Provider{
			Schema: map[string]*schema.Schema{
				"config_path": {
					Description: "The Loft config file path. Defaults to `$HOME/.loft/config.json`.",
					Type:        schema.TypeString,
					Optional:    true,
					Default:     defaultConfigPath(),
				},
				"host": {
					Description:  "The Loft instance host.",
					Type:         schema.TypeString,
					Optional:     true,
					RequiredWith: []string{"access_key"},
				},
				"insecure": {
					Description: "Allow login into an insecure Loft instance. Defaults to `false`.",
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     false,
				},
				"access_key": {
					Description:  "The Loft [access key](https://loft.sh/docs/api/access-keys).",
					Type:         schema.TypeString,
					Optional:     true,
					RequiredWith: []string{"host"},
				},
			},
			ResourcesMap: map[string]*schema.Resource{
				"loft_space":                    legacy.ResourceSpace(),
				"loft_virtual_cluster":          legacy.ResourceVirtualCluster(),
				"loft_project":                  resources.ProjectResource(),
				"loft_space_instance":           resources.SpaceInstanceResource(),
				"loft_virtual_cluster_instance": resources.VirtualClusterInstanceResource(),
				"loft_virtual_cluster_template": resources.VirtualClusterTemplateResource(),
			},
			DataSourcesMap: map[string]*schema.Resource{
				"loft_spaces":                   legacy.DataSourceSpaces(),
				"loft_space":                    legacy.DataSourceSpace(),
				"loft_virtual_cluster":          legacy.DataSourceVirtualCluster(),
				"loft_virtual_clusters":         legacy.DataSourceVirtualClusters(),
				"loft_project":                  resources.ProjectDataSource(),
				"loft_space_instance":           resources.SpaceInstanceDataSource(),
				"loft_virtual_cluster_instance": resources.VirtualClusterInstanceDataSource(),
				"loft_virtual_cluster_template": resources.VirtualClusterTemplateDataSource(),
			},
			ConfigureContextFunc: func(c context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
				var (
					loftClient client.Client
					err        error
				)

				configPath := d.Get("config_path").(string)
				if configPath != "" {
					loftClient, err = client.NewClientFromPath(configPath)
					if err != nil {
						return nil, diag.FromErr(err)
					}
				} else {
					loftClient = client.NewClient()
				}

				// Login if access key is provided
				accessKey := d.Get("access_key").(string)
				if accessKey != "" {
					host := d.Get("host").(string)
					insecure := d.Get("insecure").(bool)
					err := loftClient.LoginWithAccessKey(host, accessKey, insecure)
					if err != nil {
						return nil, diag.FromErr(err)
					}
				}

				return loftClient, nil
			},
		}
	}
}

func defaultConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(homeDir, ".loft", "config.json")
}
