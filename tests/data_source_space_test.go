package tests

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

	kubeClient, err := newKubeClient()
	if err != nil {
		t.Fatal(err)
	}

	_, adminAccessKey, configPath, err := loginUser(kubeClient, user)
	if err != nil {
		t.Fatal(err)
	}
	defer logout(t, kubeClient, adminAccessKey)

	resource.Test(t, resource.TestCase{
		CheckDestroy:      testAccSpaceCheckDestroy(kubeClient),
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSpaceCreateWithUser(configPath, user, cluster, spaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test_user", "cluster", cluster),
					resource.TestMatchResourceAttr("loft_space.test_user", "name", regexp.MustCompile(`^myspace-.*`)),
					resource.TestCheckResourceAttr("loft_space.test_user", "user", user),
					resource.TestCheckResourceAttr("loft_space.test_user", "team", ""),
				),
			},
			{
				Config: testAccResourceSpaceCreateWithUser(configPath, user, cluster, spaceName) +
					testAccDataSourceSpaceRead(cluster, spaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.loft_space.test", "cluster", cluster),
					resource.TestMatchResourceAttr("data.loft_space.test", "name", regexp.MustCompile(`^myspace-.*`)),
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
	defer logout(t, kubeClient, adminAccessKey)

	teamAccessKey, _, _, err := loginTeam(kubeClient, loftClient, cluster, team)
	if err != nil {
		t.Fatal(err)
	}
	defer logout(t, kubeClient, teamAccessKey)

	resource.Test(t, resource.TestCase{
		CheckDestroy:      testAccSpaceCheckDestroy(kubeClient),
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSpaceCreateWithTeam(configPath, team, cluster, spaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test_team", "cluster", cluster),
					resource.TestMatchResourceAttr("loft_space.test_team", "name", regexp.MustCompile(`^myspace-.*`)),
					resource.TestCheckResourceAttr("loft_space.test_team", "team", team),
					resource.TestCheckResourceAttr("loft_space.test_team", "user", ""),
				),
			},
			{
				Config: testAccResourceSpaceCreateWithTeam(configPath, team, cluster, spaceName) +
					testAccDataSourceSpaceRead(cluster, spaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.loft_space.test", "cluster", cluster),
					resource.TestMatchResourceAttr("data.loft_space.test", "name", regexp.MustCompile(`^myspace-.*`)),
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
