package provider

import (
	"fmt"
	storagev1 "github.com/loft-sh/api/v2/pkg/apis/storage/v1"
	"github.com/loft-sh/loftctl/v2/pkg/client"
	"github.com/loft-sh/loftctl/v2/pkg/kube"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"k8s.io/apiserver/pkg/storage/names"
)

func TestAccDataSourceSpace_user(t *testing.T) {
	user := "admin"
	cluster := "loft-cluster"
	spaceName := names.SimpleNameGenerator.GenerateName("myspace-")

	client, err := newKubeClient()
	if err != nil {
		t.Fatal(err)
	}

	_, adminAccessKey, configPath, err := loginUser(client, user)
	if err != nil {
		t.Fatal(err)
	}
	defer func(c kube.Interface, accessKey *storagev1.AccessKey) {
		err = logout(c, accessKey)
		if err != nil {
			fmt.Println(err)
		}
	}(client, adminAccessKey)

	resource.Test(t, resource.TestCase{
		CheckDestroy:      testAccSpaceCheckDestroy(client),
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSpaceCreate_withUser(configPath, user, cluster, spaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test_user", "cluster", cluster),
					resource.TestMatchResourceAttr("loft_space.test_user", "name", regexp.MustCompile(`^myspace\-.*`)),
					resource.TestCheckResourceAttr("loft_space.test_user", "user", user),
					resource.TestCheckResourceAttr("loft_space.test_user", "team", ""),
				),
			},
			{
				Config: testAccDataSourceSpaceRead(configPath, cluster, spaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.loft_space.test", "cluster", cluster),
					resource.TestMatchResourceAttr("data.loft_space.test", "name", regexp.MustCompile(`^myspace\-.*`)),
					resource.TestCheckResourceAttr("data.loft_space.test", "user", user),
					resource.TestCheckResourceAttr("data.loft_space.test", "team", ""),
				),
			},
		},
	})
}

func TestAccDataSourceSpace_team(t *testing.T) {
	user := "admin"
	team := "loft-admins"
	cluster := "loft-cluster"
	spaceName := names.SimpleNameGenerator.GenerateName("myspace-")

	kubeClient, err := newKubeClient()
	if err != nil {
		t.Fatal(err)
	}

	loftClient, adminAccessKey, configPath, err := loginUser(kubeClient, user)
	if err != nil {
		t.Fatal(err)
	}
	defer func(c kube.Interface, accessKey *storagev1.AccessKey) {
		err := logout(c, accessKey)
		if err != nil {
			fmt.Println(err)
		}
	}(kubeClient, adminAccessKey)

	teamAccessKey, clusterAccess, _, err := loginTeam(kubeClient, loftClient, cluster, team)
	if err != nil {
		t.Fatal(err)
	}
	defer func(c kube.Interface, accessKey *storagev1.AccessKey) {
		err := logout(c, accessKey)
		if err != nil {
			fmt.Println(err)
		}
	}(kubeClient, teamAccessKey)
	defer func(c client.Client, clusterName string, teamName string) {
		err := deleteClusterAccess(c, clusterName, teamName)
		if err != nil {
			fmt.Println(err)
		}
	}(loftClient, cluster, clusterAccess.GetName())

	resource.Test(t, resource.TestCase{
		CheckDestroy:      testAccSpaceCheckDestroy(kubeClient),
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSpaceCreate_withTeam(configPath, team, cluster, spaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test_team", "cluster", cluster),
					resource.TestMatchResourceAttr("loft_space.test_team", "name", regexp.MustCompile(`^myspace\-.*`)),
					resource.TestCheckResourceAttr("loft_space.test_team", "team", team),
					resource.TestCheckResourceAttr("loft_space.test_team", "user", ""),
				),
			},
			{
				Config: testAccDataSourceSpaceRead(configPath, cluster, spaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.loft_space.test", "cluster", cluster),
					resource.TestMatchResourceAttr("data.loft_space.test", "name", regexp.MustCompile(`^myspace\-.*`)),
					resource.TestCheckResourceAttr("data.loft_space.test", "team", team),
					resource.TestCheckResourceAttr("data.loft_space.test", "user", ""),
				),
			},
		},
	})
}

func testAccDataSourceSpaceRead(configPath string, clusterName, spaceName string) string {
	return fmt.Sprintf(`
terraform {
	required_providers {
		loft = {
			source = "registry.terraform.io/loft-sh/loft"
		}
	}
}

provider "loft" {
	config_path = "%s"
}

data "loft_space" "test" {
	cluster = "%s"
	name = "%s"
}
`,
		configPath,
		clusterName,
		spaceName,
	)
}
