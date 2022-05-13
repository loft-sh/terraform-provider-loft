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
	defer logout(t, client, adminAccessKey)

	resource.Test(t, resource.TestCase{
		CheckDestroy:      testAccSpaceCheckDestroy(client),
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSpaceCreateWithUser(configPath, user, cluster, spaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test_user", "cluster", cluster),
					resource.TestMatchResourceAttr("loft_space.test_user", "name", regexp.MustCompile(`^myspace\-.*`)),
					resource.TestCheckResourceAttr("loft_space.test_user", "user", user),
					resource.TestCheckResourceAttr("loft_space.test_user", "team", ""),
				),
			},
			{
				Config: testAccDataSourceSpaceCreateWithUser(configPath, user, cluster, spaceName) +
					testAccDataSourceSpaceRead(cluster, spaceName),
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
	defer logout(t, client, adminAccessKey)

	teamAccessKey, clusterAccess, _, err := loginTeam(client, loftClient, cluster, team)
	if err != nil {
		t.Fatal(err)
	}
	defer logout(t, client, teamAccessKey)
	defer deleteClusterAccess(t, loftClient, cluster, clusterAccess.GetName())

	resource.Test(t, resource.TestCase{
		CheckDestroy:      testAccSpaceCheckDestroy(client),
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSpaceCreateWithTeam(configPath, team, cluster, spaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test_team", "cluster", cluster),
					resource.TestMatchResourceAttr("loft_space.test_team", "name", regexp.MustCompile(`^myspace\-.*`)),
					resource.TestCheckResourceAttr("loft_space.test_team", "team", team),
					resource.TestCheckResourceAttr("loft_space.test_team", "user", ""),
				),
			},
			{
				Config: testAccDataSourceSpaceCreateWithTeam(configPath, team, cluster, spaceName) +
					testAccDataSourceSpaceRead(cluster, spaceName),
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

func testAccDataSourceSpaceRead(clusterName, spaceName string) string {
	return fmt.Sprintf(`
data "loft_space" "test" {
	cluster = "%s"
	name = "%s"
}
`,
		clusterName,
		spaceName,
	)
}
