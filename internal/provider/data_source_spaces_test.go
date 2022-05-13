package provider

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"k8s.io/apiserver/pkg/storage/names"
)

func TestAccDataSourceSpaces_all(t *testing.T) {
	rxPosNum := regexp.MustCompile("^[1-9][0-9]*$")
	user := "admin"
	team := "loft-admins"
	clusterName := "loft-cluster"
	space1Name := names.SimpleNameGenerator.GenerateName("myspace1-")
	space2Name := names.SimpleNameGenerator.GenerateName("myspace2-")
	space3Name := names.SimpleNameGenerator.GenerateName("myspace3-")
	annotation := names.SimpleNameGenerator.GenerateName("annotation-")
	space4Name := names.SimpleNameGenerator.GenerateName("objects-")
	objects := `apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config-map
data:
  foo: bar
`

	client, err := newKubeClient()
	if err != nil {
		t.Error(err)
		return
	}

	loftClient, adminAccessKey, configPath, err := loginUser(client, user)
	if err != nil {
		t.Fatal(err)
	}
	defer logout(t, client, adminAccessKey)

	teamAccessKey, clusterAccess, _, err := loginTeam(client, loftClient, clusterName, team)
	if err != nil {
		t.Fatal(err)
	}
	defer logout(t, client, teamAccessKey)
	defer deleteClusterAccess(t, loftClient, clusterName, clusterAccess.GetName())

	resource.Test(t, resource.TestCase{
		CheckDestroy:      testAccSpaceCheckDestroy(client),
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSpacesCreate(configPath, clusterName, space1Name, user, space2Name, team, space3Name, annotation, space4Name, objects),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test_user", "cluster", clusterName),
					resource.TestCheckResourceAttr("loft_space.test_user", "name", space1Name),
					resource.TestCheckResourceAttr("loft_space.test_team", "cluster", clusterName),
					resource.TestCheckResourceAttr("loft_space.test_team", "name", space2Name),
					resource.TestCheckResourceAttr("loft_space.test_annotations", "cluster", clusterName),
					resource.TestCheckResourceAttr("loft_space.test_annotations", "name", space3Name),
					resource.TestCheckResourceAttr("loft_space.test_objects", "cluster", clusterName),
					resource.TestCheckResourceAttr("loft_space.test_objects", "name", space4Name),
				),
			},
			{
				Config: testAccDataSourceSpacesCreate(configPath, clusterName, space1Name, user, space2Name, team, space3Name, annotation, space4Name, objects) +
					testAccDataSourceSpacesAll(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.loft_spaces.all", "spaces.#", rxPosNum),
					checkSpaceByName("data.loft_spaces.all", space1Name, "name", space1Name),
					checkSpaceByName("data.loft_spaces.all", space1Name, "cluster", clusterName),
					checkSpaceByName("data.loft_spaces.all", space1Name, "team", ""),
					checkSpaceByName("data.loft_spaces.all", space1Name, "user", user),
					checkSpaceByName("data.loft_spaces.all", space2Name, "name", space2Name),
					checkSpaceByName("data.loft_spaces.all", space2Name, "cluster", clusterName),
					checkSpaceByName("data.loft_spaces.all", space2Name, "team", team),
					checkSpaceByName("data.loft_spaces.all", space2Name, "user", ""),
					checkSpaceByName("data.loft_spaces.all", space3Name, "name", space3Name),
					checkSpaceByName("data.loft_spaces.all", space3Name, "cluster", clusterName),
					checkSpaceByName("data.loft_spaces.all", space3Name, "annotations.some.domain/test", annotation),
					checkSpaceByName("data.loft_spaces.all", space4Name, "name", space4Name),
					checkSpaceByName("data.loft_spaces.all", space4Name, "cluster", clusterName),
					checkSpaceByName("data.loft_spaces.all", space4Name, "user", ""),
					checkSpaceByName("data.loft_spaces.all", space4Name, "team", ""),
					checkSpaceByName("data.loft_spaces.all", space4Name, "objects", objects),
				),
			},
		},
	})
}

func testAccDataSourceSpacesCreate(
	configPath,
	clusterName,
	space1Name,
	user,
	space2Name,
	team,
	space3Name,
	annotation,
	space4Name,
	objects string,
) string {
	return fmt.Sprintf(`
terraform {
	required_providers {
		loft = {
			source = "registry.terraform.io/loft-sh/loft"
		}
	}
}

provider "loft" {
	config_path = "%[1]s"
}

resource "loft_space" "test_user" {
	name = "%[3]s"
	cluster = "%[2]s"
	user = "%[4]s"
}

resource "loft_space" "test_team" {
	name = "%[5]s"
	cluster = "%[2]s"
	team = "%[6]s"
}

resource "loft_space" "test_annotations" {
	name = "%[7]s"
	cluster = "%[2]s"
	annotations = {
		"some.domain/test" = "%[8]s"
	}
}

resource "loft_space" "test_objects" {
	name = "%[9]s"
	cluster = "%[2]s"
	objects = <<YAML
%[10]sYAML
}
`,
		configPath,
		clusterName,
		space1Name,
		user,
		space2Name,
		team,
		space3Name,
		annotation,
		space4Name,
		objects,
	)
}

func testAccDataSourceSpacesAll(clusterName string) string {
	return fmt.Sprintf(`
data "loft_spaces" "all" {
	cluster = "%s"
}
`,
		clusterName,
	)
}

func checkSpaceByName(moduleName, spaceName, key, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		spacePath := ""
		spaceNameMatch := regexp.MustCompile(`spaces\.\d+\.name`)

		primaryModule := s.RootModule().Resources[moduleName].Primary
		for key, value := range primaryModule.Attributes {
			if spaceNameMatch.MatchString(key) && value == spaceName {
				tokens := strings.Split(key, ".")
				spacePath = strings.Join([]string{tokens[0], tokens[1]}, ".")
				break
			}
		}

		if spacePath == "" {
			return fmt.Errorf("space with name %s not found", spaceName)
		}

		attrKey := strings.Join([]string{spacePath, key}, ".")
		if primaryModule.Attributes[attrKey] != value {
			return fmt.Errorf(
				"%s: Attribute '%s' didn't match %q, got %#v",
				moduleName,
				attrKey,
				value,
				primaryModule.Attributes[key])
		}

		return nil
	}
}
