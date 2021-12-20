package provider

import (
	"fmt"
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
	defer logout(client, adminAccessKey)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSpaceCreate_withUser(configPath, user, cluster, spaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
					resource.TestMatchResourceAttr("loft_space.test", "name", regexp.MustCompile(`^myspace\-.*`)),
					resource.TestCheckResourceAttr("loft_space.test", "user", user),
					resource.TestCheckResourceAttr("loft_space.test", "team", ""),
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

	client, err := newKubeClient()
	if err != nil {
		t.Fatal(err)
	}

	loftClient, adminAccessKey, configPath, err := loginUser(client, user)
	if err != nil {
		t.Fatal(err)
	}
	defer logout(client, adminAccessKey)

	teamAccessKey, clusterAccess, _, err := loginTeam(client, loftClient, cluster, team)
	if err != nil {
		t.Fatal(err)
	}
	defer logout(client, teamAccessKey)
	defer deleteClusterAccess(loftClient, cluster, clusterAccess.GetName())

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSpaceCreate_withTeam(configPath, team, cluster, spaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
					resource.TestMatchResourceAttr("loft_space.test", "name", regexp.MustCompile(`^myspace\-.*`)),
					resource.TestCheckResourceAttr("loft_space.test", "team", team),
					resource.TestCheckResourceAttr("loft_space.test", "user", ""),
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

func testAccDataSourceSpaceCreate_withUser(configPath string, user, clusterName, spaceName string) string {
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

resource "loft_space" "test" {
	name = "%s"
	cluster = "%s"
	user = "%s"
}
`,
		configPath,
		spaceName,
		clusterName,
		user,
	)
}

func testAccDataSourceSpaceCreate_withTeam(configPath string, team, clusterName, spaceName string) string {
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

resource "loft_space" "test" {
	name = "%s"
	cluster = "%s"
	team = "%s"
}
`,
		configPath,
		spaceName,
		clusterName,
		team,
	)
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
