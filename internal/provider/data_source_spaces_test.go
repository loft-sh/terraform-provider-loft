package provider

import (
	"fmt"
	storagev1 "github.com/loft-sh/api/v2/pkg/apis/storage/v1"
	"github.com/loft-sh/loftctl/v2/pkg/client"
	"github.com/loft-sh/loftctl/v2/pkg/kube"
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
	objectsName := names.SimpleNameGenerator.GenerateName("objects-")
	objects := `apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config-map
data:
  foo: bar
`

	kubeClient, err := newKubeClient()
	if err != nil {
		t.Error(err)
		return
	}

	loftClient, adminAccessKey, configPath, err := loginUser(kubeClient, user)
	if err != nil {
		t.Fatal(err)
	}
	defer func(kubeClient kube.Interface, accessKey *storagev1.AccessKey) {
		err := logout(kubeClient, accessKey)
		if err != nil {
			fmt.Println(err)
		}
	}(kubeClient, adminAccessKey)

	teamAccessKey, clusterAccess, _, err := loginTeam(kubeClient, loftClient, clusterName, team)
	if err != nil {
		t.Fatal(err)
	}
	defer func(kubeClient kube.Interface, accessKey *storagev1.AccessKey) {
		err := logout(kubeClient, accessKey)
		if err != nil {
			fmt.Println(err)
		}
	}(kubeClient, teamAccessKey)
	defer func(loftClient client.Client, clusterName string, teamName string) {
		err := deleteClusterAccess(loftClient, clusterName, teamName)
		if err != nil {
			fmt.Println(err)
		}
	}(loftClient, clusterName, clusterAccess.GetName())

	resource.Test(t, resource.TestCase{
		CheckDestroy:      testAccSpaceCheckDestroy(kubeClient),
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSpaceCreate_withUser(configPath, user, clusterName, space1Name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test_user", "cluster", clusterName),
					resource.TestCheckResourceAttr("loft_space.test_user", "name", space1Name),
				),
			},
			{
				Config: testAccDataSourceSpaceCreate_withTeam(configPath, team, clusterName, space2Name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test_team", "cluster", clusterName),
					resource.TestCheckResourceAttr("loft_space.test_team", "name", space2Name),
				),
			},
			{
				Config: testAccDataSourceSpaceCreate_withAnnotations(configPath, clusterName, space3Name, annotation),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test", "cluster", clusterName),
					resource.TestCheckResourceAttr("loft_space.test", "name", space3Name),
				),
			},
			{
				Config: testAccDataSourceSpaceCreate_withSpaceObjects(configPath, clusterName, objectsName, objects),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test_objects", "cluster", clusterName),
					resource.TestCheckResourceAttr("loft_space.test_objects", "name", objectsName),
				),
			},
			{
				Config: testAccDataSourceSpacesAll(configPath, clusterName),
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
					checkSpaceByName("data.loft_spaces.all", objectsName, "name", objectsName),
					checkSpaceByName("data.loft_spaces.all", objectsName, "cluster", clusterName),
					checkSpaceByName("data.loft_spaces.all", objectsName, "user", ""),
					checkSpaceByName("data.loft_spaces.all", objectsName, "team", ""),
					checkSpaceByName("data.loft_spaces.all", objectsName, "objects", objects),
				),
			},
		},
	})
}

func testAccDataSourceSpacesAll(configPath string, clusterName string) string {
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

data "loft_spaces" "all" {
	cluster = "%s"
}
`,
		configPath,
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
