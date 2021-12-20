package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"k8s.io/apiserver/pkg/storage/names"
)

func TestAccResourceSpace_noName(t *testing.T) {
	cluster := "loft-cluster"

	client, err := newKubeClient()
	if err != nil {
		t.Fatal(err)
	}

	_, accessKey, configPath, err := loginUser(client, "admin")
	if err != nil {
		t.Fatal(err)
	}

	defer logout(client, accessKey)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceSpaceNoName(configPath, cluster),
				ExpectError: regexp.MustCompile(`The argument "name" is required, but no definition was found.`),
			},
		},
	})
}

func TestAccResourceSpace_noCluster(t *testing.T) {
	name := names.SimpleNameGenerator.GenerateName("mycluster-")

	client, err := newKubeClient()
	if err != nil {
		t.Fatal(err)
	}

	_, accessKey, configPath, err := loginUser(client, "admin")
	if err != nil {
		t.Fatal(err)
	}

	defer logout(client, accessKey)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceSpaceNoCluster(configPath, name),
				ExpectError: regexp.MustCompile(`The argument "cluster" is required, but no definition was found.`),
			},
		},
	})
}

func TestAccResourceSpace_withGivenUser(t *testing.T) {
	name := names.SimpleNameGenerator.GenerateName("mycluster-")
	user := "admin"
	cluster := "loft-cluster"

	client, err := newKubeClient()
	if err != nil {
		t.Fatal(err)
	}

	_, accessKey, configPath, err := loginUser(client, user)
	if err != nil {
		t.Fatal(err)
	}

	defer logout(client, accessKey)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSpaceUser(configPath, name, cluster, user),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test", "user", user),
				),
			},
		},
	})
}

func TestAccResourceSpace_withGivenTeam(t *testing.T) {
	name := names.SimpleNameGenerator.GenerateName("mycluster-")
	cluster := "loft-cluster"
	user := "admin"
	team := "loft-admins"

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
				Config: testAccResourceSpaceTeam(configPath, name, cluster, team),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
				),
			},
		},
	})
}

func testAccResourceSpaceNoName(configPath, clusterName string) string {
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
	cluster = "%s"
}
`,
		configPath,
		clusterName,
	)
}

func testAccResourceSpaceNoCluster(configPath, spaceName string) string {
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
}
`,
		configPath,
		spaceName,
	)
}

func testAccResourceSpaceUser(configPath, spaceName, clusterName, user string) string {
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

func testAccResourceSpaceTeam(configPath, spaceName, clusterName, team string) string {
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
