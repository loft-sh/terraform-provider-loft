package provider

import (
	"context"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/loft-sh/loftctl/v2/pkg/client"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			DataSourcesMap: map[string]*schema.Resource{
				"loft_spaces":           DataSourceSpaces(),
				"loft_space":            DataSourceSpace(),
				"loft_virtual_cluster":  DataSourceVirtualCluster(),
				"loft_virtual_clusters": DataSourceVirtualClusters(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"loft_space":           ResourceSpace(),
				"loft_virtual_cluster": ResourceVirtualCluster(),
			},
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
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

type apiClient struct {
	LoftClient client.Client
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(c context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
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

		apiClient := &apiClient{
			LoftClient: loftClient,
		}

		return apiClient, nil
	}
}

func defaultConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(homeDir, ".loft", "config.json")
}
