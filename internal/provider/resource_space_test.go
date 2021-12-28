package provider

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	v1 "github.com/loft-sh/agentapi/v2/pkg/apis/loft/cluster/v1"
	"github.com/loft-sh/loftctl/v2/pkg/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
				Config: testAccDataSourceSpaceCreate_withUser(configPath, user, cluster, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test_user", "name", name),
					resource.TestCheckResourceAttr("loft_space.test_user", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test_user", "user", user),
					resource.TestCheckResourceAttr("loft_space.test_user", "team", ""),
					checkSpace(configPath, cluster, name, hasUser(user)),
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
				Config: testAccDataSourceSpaceCreate_withTeam(configPath, team, cluster, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test_team", "name", name),
					resource.TestCheckResourceAttr("loft_space.test_team", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test_team", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test_team", "team", team),
					checkSpace(configPath, cluster, name, hasTeam(team)),
				),
			},
		},
	})
}

func TestAccResourceSpace_withAnnotations(t *testing.T) {
	name := names.SimpleNameGenerator.GenerateName("mycluster-")
	annotation := names.SimpleNameGenerator.GenerateName("annotation-")
	cluster := "loft-cluster"
	user := "admin"

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
				Config: testAccDataSourceSpaceCreate_withAnnotations(configPath, cluster, name, annotation),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test_annotations", "name", name),
					resource.TestCheckResourceAttr("loft_space.test_annotations", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test_annotations", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test_annotations", "team", ""),
					resource.TestCheckResourceAttr("loft_space.test_annotations", "annotations.loft.sh/test", annotation),
					checkSpace(configPath, cluster, name, hasAnnotation("loft.sh/test", annotation)),
				),
			},
		},
	})
}

func TestAccResourceSpace_withLabels(t *testing.T) {
	name := names.SimpleNameGenerator.GenerateName("mycluster-")
	label := names.SimpleNameGenerator.GenerateName("annotation-")
	cluster := "loft-cluster"
	user := "admin"

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
				Config: testAccDataSourceSpaceCreate_withLabels(configPath, cluster, name, label),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test_labels", "name", name),
					resource.TestCheckResourceAttr("loft_space.test_labels", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test_labels", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test_labels", "team", ""),
					resource.TestCheckResourceAttr("loft_space.test_labels", "labels.loft.sh/test", label),
					checkSpace(configPath, cluster, name, hasLabel("loft.sh/test", label)),
				),
			},
		},
	})
}

func TestAccResourceSpace_withSleepAfter(t *testing.T) {
	rxPosNum := regexp.MustCompile("^[1-9][0-9]*$")
	name := names.SimpleNameGenerator.GenerateName("mycluster-")
	cluster := "loft-cluster"
	user := "admin"
	sleepAfter := 60

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
				Config: testAccDataSourceSpaceCreate_withSleepAfter(configPath, cluster, name, sleepAfter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test_sleep_after", "name", name),
					resource.TestCheckResourceAttr("loft_space.test_sleep_after", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test_sleep_after", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test_sleep_after", "team", ""),
					resource.TestMatchResourceAttr("loft_space.test_sleep_after", "sleep_after", rxPosNum),
					checkSpace(configPath, cluster, name, hasAnnotation(v1.SleepModeSleepAfterAnnotation, strconv.Itoa(sleepAfter))),
				),
			},
		},
	})
}

func TestAccResourceSpace_withDeleteAfter(t *testing.T) {
	rxPosNum := regexp.MustCompile("^[1-9][0-9]*$")
	name := names.SimpleNameGenerator.GenerateName("mycluster-")
	cluster := "loft-cluster"
	user := "admin"
	deleteAfter := 60

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
				Config: testAccDataSourceSpaceCreate_withDeleteAfter(configPath, cluster, name, deleteAfter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test_delete_after", "name", name),
					resource.TestCheckResourceAttr("loft_space.test_delete_after", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test_delete_after", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test_delete_after", "team", ""),
					resource.TestMatchResourceAttr("loft_space.test_delete_after", "delete_after", rxPosNum),
					checkSpace(configPath, cluster, name, hasAnnotation(v1.SleepModeDeleteAfterAnnotation, strconv.Itoa(deleteAfter))),
				),
			},
		},
	})
}

func TestAccResourceSpace_withSleepSchedule(t *testing.T) {
	name := names.SimpleNameGenerator.GenerateName("mycluster-")
	cluster := "loft-cluster"
	user := "admin"
	sleepSchedule := "0 0 * * *"

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
				Config: testAccDataSourceSpaceCreate_withScheduledSleep(configPath, cluster, name, sleepSchedule),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test_sleep_schedule", "name", name),
					resource.TestCheckResourceAttr("loft_space.test_sleep_schedule", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test_sleep_schedule", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test_sleep_schedule", "team", ""),
					resource.TestCheckResourceAttr("loft_space.test_sleep_schedule", "sleep_schedule", sleepSchedule),
					checkSpace(configPath, cluster, name, hasAnnotation(v1.SleepModeSleepScheduleAnnotation, sleepSchedule)),
				),
			},
		},
	})
}

func TestAccResourceSpace_withWakeSchedule(t *testing.T) {
	name := names.SimpleNameGenerator.GenerateName("mycluster-")
	cluster := "loft-cluster"
	user := "admin"
	wakeSchedule := "0 0 * * *"

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
				Config: testAccDataSourceSpaceCreate_withScheduledWakeup(configPath, cluster, name, wakeSchedule),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test_wakeup_schedule", "name", name),
					resource.TestCheckResourceAttr("loft_space.test_wakeup_schedule", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test_wakeup_schedule", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test_wakeup_schedule", "team", ""),
					resource.TestCheckResourceAttr("loft_space.test_wakeup_schedule", "wakeup_schedule", wakeSchedule),
					checkSpace(configPath, cluster, name, hasAnnotation(v1.SleepModeWakeupScheduleAnnotation, wakeSchedule)),
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

func checkSpace(configPath, clusterName, spaceName string, pred func(space *v1.Space) error) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		apiClient, err := client.NewClientFromPath(configPath)
		if err != nil {
			return err
		}

		clusterClient, err := apiClient.Cluster(clusterName)
		if err != nil {
			return err
		}

		space, err := clusterClient.Agent().ClusterV1().Spaces().Get(context.TODO(), spaceName, metav1.GetOptions{})
		if err != nil {
			return err
		}

		return pred(space)
	}
}

func hasAnnotation(annotation, value string) func(space *v1.Space) error {
	return func(space *v1.Space) error {
		if space.GetAnnotations()[annotation] != value {
			return fmt.Errorf(
				"%s: Annotation '%s' didn't match %q, got %#v",
				space.GetName(),
				annotation,
				value,
				space.GetLabels()[annotation])
		}
		return nil
	}
}

func hasLabel(label, value string) func(space *v1.Space) error {
	return func(space *v1.Space) error {
		if space.GetLabels()[label] != value {
			return fmt.Errorf(
				"%s: Label '%s' didn't match %q, got %#v",
				space.GetName(),
				label,
				value,
				space.GetLabels()[label])
		}
		return nil
	}
}

func hasUser(user string) func(space *v1.Space) error {
	return func(space *v1.Space) error {
		if space.Spec.User != user {
			return fmt.Errorf(
				"%s: User didn't match %q, got %#v",
				space.GetName(),
				user,
				space.Spec.User)
		}
		return nil
	}
}

func hasTeam(team string) func(space *v1.Space) error {
	return func(space *v1.Space) error {
		if space.Spec.Team != team {
			return fmt.Errorf(
				"%s: Team didn't match %q, got %#v",
				space.GetName(),
				team,
				space.Spec.Team)
		}
		return nil
	}
}
