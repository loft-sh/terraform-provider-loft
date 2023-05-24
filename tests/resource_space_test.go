package tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/loft-sh/terraform-provider-loft/internal/provider"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	v1 "github.com/loft-sh/agentapi/v3/pkg/apis/loft/cluster/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apiserver/pkg/storage/names"
)

func TestAccResourceSpace_noNameOrGenerateName(t *testing.T) {
	cluster := "loft-cluster"

	kubeClient, err := newKubeClient()
	if err != nil {
		t.Fatal(err)
	}

	_, accessKey, configPath, err := loginUser(kubeClient, "admin")
	if err != nil {
		t.Fatal(err)
	}

	defer logout(t, kubeClient, accessKey)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceSpaceNoName(configPath, cluster),
				ExpectError: regexp.MustCompile(`Required value: name or generateName is required`),
			},
		},
	})
}

func TestAccResourceSpace_noCluster(t *testing.T) {
	name := names.SimpleNameGenerator.GenerateName("mycluster-")

	kubeClient, err := newKubeClient()
	if err != nil {
		t.Fatal(err)
	}

	_, accessKey, configPath, err := loginUser(kubeClient, "admin")
	if err != nil {
		t.Fatal(err)
	}

	defer logout(t, kubeClient, accessKey)

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

	kubeClient, err := newKubeClient()
	if err != nil {
		t.Fatal(err)
	}

	_, accessKey, configPath, err := loginUser(kubeClient, user)
	if err != nil {
		t.Fatal(err)
	}
	defer logout(t, kubeClient, accessKey)

	resource.Test(t, resource.TestCase{
		CheckDestroy:      testAccSpaceCheckDestroy(kubeClient),
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSpaceCreateWithUser(configPath, user, cluster, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test_user", "name", name),
					resource.TestCheckResourceAttr("loft_space.test_user", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test_user", "user", user),
					resource.TestCheckResourceAttr("loft_space.test_user", "team", ""),
					checkSpace(configPath, cluster, name, spaceHasUser(user)),
				),
			},
			{
				ResourceName:      "loft_space.test_user",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccResourceSpaceCreateWithUser(configPath, user2, cluster, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test_user", "name", name),
					resource.TestCheckResourceAttr("loft_space.test_user", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test_user", "user", user2),
					resource.TestCheckResourceAttr("loft_space.test_user", "team", ""),
					checkSpace(configPath, cluster, name, spaceHasUser(user2)),
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
				Config: testAccResourceSpaceCreateWithTeam(configPath, team, cluster, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test_team", "name", name),
					resource.TestCheckResourceAttr("loft_space.test_team", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test_team", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test_team", "team", team),
					checkSpace(configPath, cluster, name, spaceHasTeam(team)),
				),
			},
			{
				ResourceName:      "loft_space.test_team",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccResourceSpaceCreateWithTeam(configPath, team2, cluster, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test_team", "name", name),
					resource.TestCheckResourceAttr("loft_space.test_team", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test_team", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test_team", "team", team2),
					checkSpace(configPath, cluster, name, spaceHasTeam(team2)),
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
				Config: testAccResourceSpaceCreateWithAnnotations(configPath, cluster, name, annotation),
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
				Config: testAccResourceSpaceCreateWithAnnotations(configPath, cluster, name, annotation2),
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
				Config: testAccResourceSpaceCreateWithoutUserOrTeam(configPath, cluster, name),
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
				Config: testAccResourceSpaceCreateWithLabels(configPath, cluster, name, label),
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
				Config: testAccResourceSpaceCreateWithLabels(configPath, cluster, name, label2),
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
				Config: testAccResourceSpaceCreateWithoutUserOrTeam(configPath, cluster, name),
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

func TestAccResourceSpace_withInvalidSleepAfter(t *testing.T) {
	name := names.SimpleNameGenerator.GenerateName("mycluster-")
	cluster := "loft-cluster"

	kubeClient, err := newKubeClient()
	if err != nil {
		t.Fatal(err)
	}

	_, accessKey, configPath, err := loginUser(kubeClient, "admin")
	if err != nil {
		t.Fatal(err)
	}

	defer logout(t, kubeClient, accessKey)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceSpaceCreateWithSleepAfter(configPath, cluster, name, "oops"),
				ExpectError: regexp.MustCompile(`time: invalid duration "oops"`),
			},
		},
	})
}

func TestAccResourceSpace_withSleepAfter(t *testing.T) {
	name := names.SimpleNameGenerator.GenerateName("mycluster-")
	cluster := "loft-cluster"
	user := "admin"
	sleepAfter := "1m"
	sleepAfterSeconds, err := toSecondsString(sleepAfter)
	if err != nil {
		t.Fatal(err)
	}
	sleepAfter2 := "120s"
	sleepAfter2Seconds, err := toSecondsString(sleepAfter2)
	if err != nil {
		t.Fatal(err)
	}

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
				Config: testAccResourceSpaceCreateWithSleepAfter(configPath, cluster, name, sleepAfter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test", "name", name),
					resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test", "team", ""),
					resource.TestCheckResourceAttr("loft_space.test", "sleep_after", sleepAfterSeconds),
					checkSpace(configPath, cluster, name, hasAnnotation(v1.SleepModeSleepAfterAnnotation, sleepAfterSeconds)),
				),
			},
			{
				ResourceName:      "loft_space.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccResourceSpaceCreateWithSleepAfter(configPath, cluster, name, sleepAfter2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test", "name", name),
					resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test", "team", ""),
					resource.TestCheckResourceAttr("loft_space.test", "sleep_after", sleepAfter2Seconds),
					checkSpace(configPath, cluster, name, hasAnnotation(v1.SleepModeSleepAfterAnnotation, sleepAfter2Seconds)),
				),
			},
			{
				ResourceName:      "loft_space.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccResourceSpaceCreateWithoutUserOrTeam(configPath, cluster, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test", "name", name),
					resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test", "team", ""),
					resource.TestCheckResourceAttr("loft_space.test", "sleep_after", ""),
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

func TestAccResourceSpace_withInvalidDeleteAfter(t *testing.T) {
	name := names.SimpleNameGenerator.GenerateName("mycluster-")
	cluster := "loft-cluster"

	kubeClient, err := newKubeClient()
	if err != nil {
		t.Fatal(err)
	}

	_, accessKey, configPath, err := loginUser(kubeClient, "admin")
	if err != nil {
		t.Fatal(err)
	}

	defer logout(t, kubeClient, accessKey)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceSpaceCreateWithDeleteAfter(configPath, cluster, name, "oops"),
				ExpectError: regexp.MustCompile(`time: invalid duration "oops"`),
			},
		},
	})
}

func TestAccResourceSpace_withDeleteAfter(t *testing.T) {
	name := names.SimpleNameGenerator.GenerateName("mycluster-")
	cluster := "loft-cluster"
	user := "admin"
	deleteAfter := "1m0s"
	deleteAfterSeconds, err := toSecondsString(deleteAfter)
	if err != nil {
		t.Fatal(err)
	}

	deleteAfter2 := "2m0s"
	deleteAfter2Seconds, err := toSecondsString(deleteAfter2)
	if err != nil {
		t.Fatal(err)
	}

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
				Config: testAccResourceSpaceCreateWithDeleteAfter(configPath, cluster, name, deleteAfter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test", "name", name),
					resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test", "team", ""),
					resource.TestCheckResourceAttr("loft_space.test", "delete_after", deleteAfterSeconds),
					checkSpace(configPath, cluster, name, hasAnnotation(v1.SleepModeDeleteAfterAnnotation, deleteAfterSeconds)),
				),
			},
			{
				ResourceName:      "loft_space.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccResourceSpaceCreateWithDeleteAfter(configPath, cluster, name, deleteAfter2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test", "name", name),
					resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test", "team", ""),
					resource.TestCheckResourceAttr("loft_space.test", "delete_after", deleteAfter2Seconds),
					checkSpace(configPath, cluster, name, hasAnnotation(v1.SleepModeDeleteAfterAnnotation, deleteAfter2Seconds)),
				),
			},
			{
				ResourceName:      "loft_space.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccResourceSpaceCreateWithoutUserOrTeam(configPath, cluster, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test", "name", name),
					resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test", "team", ""),
					resource.TestCheckResourceAttr("loft_space.test", "delete_after", ""),
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
				Config: testAccResourceSpaceCreateWithScheduledSleep(configPath, cluster, name, sleepSchedule),
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
				Config: testAccResourceSpaceCreateWithScheduledSleep(configPath, cluster, name, sleepSchedule2),
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
				Config: testAccResourceSpaceCreateWithoutUserOrTeam(configPath, cluster, name),
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
				Config: testAccResourceSpaceCreateWithScheduledWakeup(configPath, cluster, name, wakeSchedule),
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
				Config: testAccResourceSpaceCreateWithScheduledWakeup(configPath, cluster, name, wakeSchedule2),
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
				Config: testAccResourceSpaceCreateWithoutUserOrTeam(configPath, cluster, name),
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
				Config: testAccResourceSpaceCreateWithSpaceConstraints(configPath, cluster, name, spaceConstraints),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test", "name", name),
					resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test", "team", ""),
					resource.TestCheckResourceAttr("loft_space.test", "space_constraints", spaceConstraints),
					checkSpace(configPath, cluster, name, hasLabel(provider.SpaceLabelSpaceConstraints, spaceConstraints)),
				),
			},
			{
				ResourceName:            "loft_space.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"space_constraints"},
			},
			{
				Config: testAccResourceSpaceCreateWithSpaceConstraints(configPath, cluster, name, spaceConstraints2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test", "name", name),
					resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test", "team", ""),
					resource.TestCheckResourceAttr("loft_space.test", "space_constraints", spaceConstraints2),
					checkSpace(configPath, cluster, name, hasLabel(provider.SpaceLabelSpaceConstraints, spaceConstraints2)),
				),
			},
			{
				ResourceName:            "loft_space.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"space_constraints"},
			},
			{
				Config: testAccResourceSpaceCreateWithoutUserOrTeam(configPath, cluster, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space.test", "name", name),
					resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test", "team", ""),
					resource.TestCheckResourceAttr("loft_space.test", "space_constraints", ""),
					checkSpace(configPath, cluster, name, noLabel(provider.SpaceLabelSpaceConstraints)),
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
				Config: testAccResourceSpaceCreateWithSpaceObjects(configPath, cluster, name, objects1),
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
				Config: testAccResourceSpaceCreateWithSpaceObjects(configPath, cluster, name, objects2),
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
				Config: testAccResourceSpaceCreateWithSpaceObjects(configPath, cluster, name, ""),
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

func TestAccResourceSpace_withGenerateName(t *testing.T) {
	prefix := "test-space-"
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
				Config: testAccResourceSpaceCreateWithGenerateName(configPath, cluster, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("loft_space.test_generate_name", "name", regexp.MustCompile(prefix)),
					resource.TestCheckResourceAttr("loft_space.test_generate_name", "generate_name", prefix),
					resource.TestCheckResourceAttr("loft_space.test_generate_name", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_space.test_generate_name", "user", ""),
					resource.TestCheckResourceAttr("loft_space.test_generate_name", "team", ""),
					resource.TestCheckResourceAttr("loft_space.test_generate_name", "objects", ""),
				),
			},
			{
				ResourceName:      "loft_space.test_generate_name",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceSpace_withNameAndGenerateName(t *testing.T) {
	cluster := "loft-cluster"
	name := names.SimpleNameGenerator.GenerateName("mycluster-")
	prefix := "mycluster-"

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
				Config:      testAccResourceSpaceCreateWithNameAndGenerateName(configPath, cluster, name, prefix),
				ExpectError: regexp.MustCompile(`"generate_name": conflicts with name`),
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
