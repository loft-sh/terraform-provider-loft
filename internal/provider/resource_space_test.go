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
	corev1 "k8s.io/api/core/v1"
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

	defer logout(t, client, accessKey)

	resource.Test(t, resource.TestCase{
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

	defer logout(t, client, accessKey)

	resource.Test(t, resource.TestCase{
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
	user2 := "admin2"
	cluster := "loft-cluster"

	client, err := newKubeClient()
	if err != nil {
		t.Fatal(err)
	}

	_, accessKey, configPath, err := loginUser(client, user)
	if err != nil {
		t.Fatal(err)
	}
	defer logout(t, client, accessKey)

	resource.Test(t, resource.TestCase{
		CheckDestroy:      testAccSpaceCheckDestroy(client),
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSpaceCreateWithUser(configPath, user, cluster, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test_user", "name", name),
					resource.TestCheckResourceAttr("loft_space.test_user", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test_user", "user", user),
					resource.TestCheckResourceAttr("loft_space.test_user", "team", ""),
					checkSpace(configPath, cluster, name, hasUser(user)),
				),
			},
			{
				ResourceName:      "loft_space.test_user",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataSourceSpaceCreateWithUser(configPath, user2, cluster, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test_user", "name", name),
					resource.TestCheckResourceAttr("loft_space.test_user", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test_user", "user", user2),
					resource.TestCheckResourceAttr("loft_space.test_user", "team", ""),
					checkSpace(configPath, cluster, name, hasUser(user2)),
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
	team2 := "loft-admins2"

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
				Config: testAccDataSourceSpaceCreateWithTeam(configPath, team, cluster, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test_team", "name", name),
					resource.TestCheckResourceAttr("loft_space.test_team", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test_team", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test_team", "team", team),
					checkSpace(configPath, cluster, name, hasTeam(team)),
				),
			},
			{
				ResourceName:      "loft_space.test_team",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataSourceSpaceCreateWithTeam(configPath, team2, cluster, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test_team", "name", name),
					resource.TestCheckResourceAttr("loft_space.test_team", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test_team", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test_team", "team", team2),
					checkSpace(configPath, cluster, name, hasTeam(team2)),
				),
			},
		},
	})
}

func TestAccResourceSpace_withAnnotations(t *testing.T) {
	name := names.SimpleNameGenerator.GenerateName("mycluster-")
	annotation := names.SimpleNameGenerator.GenerateName("annotation-")
	annotation2 := names.SimpleNameGenerator.GenerateName("annotation-")
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
	defer logout(t, client, adminAccessKey)

	resource.Test(t, resource.TestCase{
		CheckDestroy:      testAccSpaceCheckDestroy(client),
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSpaceCreateWithAnnotations(configPath, cluster, name, annotation),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test", "name", name),
					resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test", "team", ""),
					resource.TestCheckResourceAttr("loft_space.test", "annotations.some.domain/test", annotation),
					checkSpace(configPath, cluster, name, hasAnnotation("some.domain/test", annotation)),
				),
			},
			{
				ResourceName:      "loft_space.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataSourceSpaceCreateWithAnnotations(configPath, cluster, name, annotation2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test", "name", name),
					resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test", "team", ""),
					resource.TestCheckResourceAttr("loft_space.test", "annotations.some.domain/test", annotation2),
					checkSpace(configPath, cluster, name, hasAnnotation("some.domain/test", annotation2)),
				),
			},
			{
				ResourceName:      "loft_space.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataSourceSpaceCreateWithoutUserOrTeam(configPath, cluster, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test", "name", name),
					resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test", "team", ""),
					checkSpace(configPath, cluster, name, noAnnotation("some.domain/test")),
				),
			},
			{
				ResourceName:      "loft_space.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceSpace_withLabels(t *testing.T) {
	name := names.SimpleNameGenerator.GenerateName("mycluster-")
	label := names.SimpleNameGenerator.GenerateName("annotation-")
	label2 := names.SimpleNameGenerator.GenerateName("annotation-")
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
	defer logout(t, client, adminAccessKey)

	resource.Test(t, resource.TestCase{
		CheckDestroy:      testAccSpaceCheckDestroy(client),
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSpaceCreateWithLabels(configPath, cluster, name, label),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test", "name", name),
					resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test", "team", ""),
					resource.TestCheckResourceAttr("loft_space.test", "labels.some.domain/test", label),
					checkSpace(configPath, cluster, name, hasLabel("some.domain/test", label)),
					checkSpace(configPath, cluster, name, hasLabel(corev1.LabelMetadataName, name)),
				),
			},
			{
				ResourceName:      "loft_space.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataSourceSpaceCreateWithLabels(configPath, cluster, name, label2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test", "name", name),
					resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test", "team", ""),
					resource.TestCheckResourceAttr("loft_space.test", "labels.some.domain/test", label2),
					checkSpace(configPath, cluster, name, hasLabel("some.domain/test", label2)),
					checkSpace(configPath, cluster, name, hasLabel(corev1.LabelMetadataName, name)),
				),
			},
			{
				ResourceName:      "loft_space.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataSourceSpaceCreateWithoutUserOrTeam(configPath, cluster, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test", "name", name),
					resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test", "team", ""),
					checkSpace(configPath, cluster, name, noLabel("some.domain/test")),
					checkSpace(configPath, cluster, name, hasLabel(corev1.LabelMetadataName, name)),
				),
			},
			{
				ResourceName:      "loft_space.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceSpace_withSleepAfter(t *testing.T) {
	name := names.SimpleNameGenerator.GenerateName("mycluster-")
	cluster := "loft-cluster"
	user := "admin"
	sleepAfter := 60
	sleepAfter2 := 120

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
				Config: testAccDataSourceSpaceCreateWithSleepAfter(configPath, cluster, name, sleepAfter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test", "name", name),
					resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test", "team", ""),
					resource.TestCheckResourceAttr("loft_space.test", "sleep_after", strconv.Itoa(sleepAfter)),
					checkSpace(configPath, cluster, name, hasAnnotation(v1.SleepModeSleepAfterAnnotation, strconv.Itoa(sleepAfter))),
				),
			},
			{
				ResourceName:      "loft_space.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataSourceSpaceCreateWithSleepAfter(configPath, cluster, name, sleepAfter2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test", "name", name),
					resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test", "team", ""),
					resource.TestCheckResourceAttr("loft_space.test", "sleep_after", strconv.Itoa(sleepAfter2)),
					checkSpace(configPath, cluster, name, hasAnnotation(v1.SleepModeSleepAfterAnnotation, strconv.Itoa(sleepAfter2))),
				),
			},
			{
				ResourceName:      "loft_space.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataSourceSpaceCreateWithoutUserOrTeam(configPath, cluster, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test", "name", name),
					resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test", "team", ""),
					resource.TestCheckResourceAttr("loft_space.test", "sleep_after", "0"),
					checkSpace(configPath, cluster, name, noAnnotation(v1.SleepModeSleepAfterAnnotation)),
				),
			},
			{
				ResourceName:            "loft_space.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"sleep_after"},
			},
		},
	})
}

func TestAccResourceSpace_withDeleteAfter(t *testing.T) {
	name := names.SimpleNameGenerator.GenerateName("mycluster-")
	cluster := "loft-cluster"
	user := "admin"
	deleteAfter := 60
	deleteAfter2 := 120

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
				Config: testAccDataSourceSpaceCreateWithDeleteAfter(configPath, cluster, name, deleteAfter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test", "name", name),
					resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test", "team", ""),
					resource.TestCheckResourceAttr("loft_space.test", "delete_after", strconv.Itoa(deleteAfter)),
					checkSpace(configPath, cluster, name, hasAnnotation(v1.SleepModeDeleteAfterAnnotation, strconv.Itoa(deleteAfter))),
				),
			},
			{
				ResourceName:      "loft_space.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataSourceSpaceCreateWithDeleteAfter(configPath, cluster, name, deleteAfter2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test", "name", name),
					resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test", "team", ""),
					resource.TestCheckResourceAttr("loft_space.test", "delete_after", strconv.Itoa(deleteAfter2)),
					checkSpace(configPath, cluster, name, hasAnnotation(v1.SleepModeDeleteAfterAnnotation, strconv.Itoa(deleteAfter2))),
				),
			},
			{
				ResourceName:      "loft_space.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataSourceSpaceCreateWithoutUserOrTeam(configPath, cluster, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test", "name", name),
					resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test", "team", ""),
					resource.TestCheckResourceAttr("loft_space.test", "delete_after", "0"),
					checkSpace(configPath, cluster, name, noAnnotation(v1.SleepModeDeleteAfterAnnotation)),
				),
			},
			{
				ResourceName:            "loft_space.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_after"},
			},
		},
	})
}

func TestAccResourceSpace_withSleepSchedule(t *testing.T) {
	name := names.SimpleNameGenerator.GenerateName("mycluster-")
	cluster := "loft-cluster"
	user := "admin"
	sleepSchedule := "0 0 * * *"
	sleepSchedule2 := "30 6 * * *"

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
				Config: testAccDataSourceSpaceCreateWithScheduledSleep(configPath, cluster, name, sleepSchedule),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test", "name", name),
					resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test", "team", ""),
					resource.TestCheckResourceAttr("loft_space.test", "sleep_schedule", sleepSchedule),
					checkSpace(configPath, cluster, name, hasAnnotation(v1.SleepModeSleepScheduleAnnotation, sleepSchedule)),
				),
			},
			{
				ResourceName:      "loft_space.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataSourceSpaceCreateWithScheduledSleep(configPath, cluster, name, sleepSchedule2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test", "name", name),
					resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test", "team", ""),
					resource.TestCheckResourceAttr("loft_space.test", "sleep_schedule", sleepSchedule2),
					checkSpace(configPath, cluster, name, hasAnnotation(v1.SleepModeSleepScheduleAnnotation, sleepSchedule2)),
				),
			},
			{
				ResourceName:      "loft_space.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataSourceSpaceCreateWithoutUserOrTeam(configPath, cluster, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test", "name", name),
					resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test", "team", ""),
					resource.TestCheckResourceAttr("loft_space.test", "sleep_schedule", ""),
					checkSpace(configPath, cluster, name, noAnnotation(v1.SleepModeSleepScheduleAnnotation)),
				),
			},
			{
				ResourceName:            "loft_space.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"sleep_schedule"},
			},
		},
	})
}

func TestAccResourceSpace_withWakeupSchedule(t *testing.T) {
	name := names.SimpleNameGenerator.GenerateName("mycluster-")
	cluster := "loft-cluster"
	user := "admin"
	wakeSchedule := "0 0 * * *"
	wakeSchedule2 := "30 18 * * *"

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
				Config: testAccDataSourceSpaceCreateWithScheduledWakeup(configPath, cluster, name, wakeSchedule),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test", "name", name),
					resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test", "team", ""),
					resource.TestCheckResourceAttr("loft_space.test", "wakeup_schedule", wakeSchedule),
					checkSpace(configPath, cluster, name, hasAnnotation(v1.SleepModeWakeupScheduleAnnotation, wakeSchedule)),
				),
			},
			{
				ResourceName:      "loft_space.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataSourceSpaceCreateWithScheduledWakeup(configPath, cluster, name, wakeSchedule2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test", "name", name),
					resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test", "team", ""),
					resource.TestCheckResourceAttr("loft_space.test", "wakeup_schedule", wakeSchedule2),
					checkSpace(configPath, cluster, name, hasAnnotation(v1.SleepModeWakeupScheduleAnnotation, wakeSchedule2)),
				),
			},
			{
				ResourceName:      "loft_space.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataSourceSpaceCreateWithoutUserOrTeam(configPath, cluster, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test", "name", name),
					resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test", "team", ""),
					resource.TestCheckResourceAttr("loft_space.test", "wakeup_schedule", ""),
					checkSpace(configPath, cluster, name, noAnnotation(v1.SleepModeWakeupScheduleAnnotation)),
				),
			},
			{
				ResourceName:            "loft_space.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"wakeup_schedule"},
			},
		},
	})
}

func TestAccResourceSpace_withSpaceConstraints(t *testing.T) {
	name := names.SimpleNameGenerator.GenerateName("mycluster-")
	cluster := "loft-cluster"
	user := "admin"
	spaceConstraints := "default"
	spaceConstraints2 := "isolated"

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
				Config: testAccDataSourceSpaceCreateWithSpaceConstraints(configPath, cluster, name, spaceConstraints),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test", "name", name),
					resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test", "team", ""),
					resource.TestCheckResourceAttr("loft_space.test", "space_constraints", spaceConstraints),
					checkSpace(configPath, cluster, name, hasLabel(SpaceLabelSpaceConstraints, spaceConstraints)),
				),
			},
			{
				ResourceName:            "loft_space.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"space_constraints"},
			},
			{
				Config: testAccDataSourceSpaceCreateWithSpaceConstraints(configPath, cluster, name, spaceConstraints2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test", "name", name),
					resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test", "team", ""),
					resource.TestCheckResourceAttr("loft_space.test", "space_constraints", spaceConstraints2),
					checkSpace(configPath, cluster, name, hasLabel(SpaceLabelSpaceConstraints, spaceConstraints2)),
				),
			},
			{
				ResourceName:            "loft_space.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"space_constraints"},
			},
			{
				Config: testAccDataSourceSpaceCreateWithoutUserOrTeam(configPath, cluster, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test", "name", name),
					resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test", "team", ""),
					resource.TestCheckResourceAttr("loft_space.test", "space_constraints", ""),
					checkSpace(configPath, cluster, name, noLabel(SpaceLabelSpaceConstraints)),
				),
			},
			{
				ResourceName:            "loft_space.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"space_constraints"},
			},
		},
	})
}

func TestAccResourceSpace_withSpaceObjects(t *testing.T) {
	name := names.SimpleNameGenerator.GenerateName("mycluster-")
	cluster := "loft-cluster"
	user := "admin"
	objects1 := `apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config-map
data:
  foo: bar
`
	objects2 := `apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config-map
data:
  foo: bar
  hello: world
`

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
				Config: testAccDataSourceSpaceCreateWithSpaceObjects(configPath, cluster, name, objects1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test_objects", "name", name),
					resource.TestCheckResourceAttr("loft_space.test_objects", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test_objects", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test_objects", "team", ""),
					resource.TestCheckResourceAttr("loft_space.test_objects", "objects", objects1),
				),
			},
			{
				ResourceName:      "loft_space.test_objects",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataSourceSpaceCreateWithSpaceObjects(configPath, cluster, name, objects2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test_objects", "name", name),
					resource.TestCheckResourceAttr("loft_space.test_objects", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test_objects", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test_objects", "team", ""),
					resource.TestCheckResourceAttr("loft_space.test_objects", "objects", objects2),
				),
			},
			{
				ResourceName:      "loft_space.test_objects",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataSourceSpaceCreateWithSpaceObjects(configPath, cluster, name, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test_objects", "name", name),
					resource.TestCheckResourceAttr("loft_space.test_objects", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test_objects", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test_objects", "team", ""),
					resource.TestCheckResourceAttr("loft_space.test_objects", "objects", ""),
				),
			},
			{
				ResourceName:      "loft_space.test_objects",
				ImportState:       true,
				ImportStateVerify: true,
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

func noAnnotation(annotation string) func(space *v1.Space) error {
	return func(space *v1.Space) error {
		if space.GetAnnotations()[annotation] != "" {
			return fmt.Errorf(
				"%s: Annotation '%s' should not be present",
				space.GetName(),
				annotation,
			)
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

func noLabel(label string) func(space *v1.Space) error {
	return func(space *v1.Space) error {
		if space.GetAnnotations()[label] != "" {
			return fmt.Errorf(
				"%s: Label '%s' should not be present",
				space.GetName(),
				label,
			)
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
