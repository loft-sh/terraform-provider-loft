package tests

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	managementv1 "github.com/loft-sh/api/v2/pkg/apis/management/v1"
	clientpkg "github.com/loft-sh/loftctl/v2/pkg/client"
	"github.com/loft-sh/loftctl/v2/pkg/client/naming"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/storage/names"
	"regexp"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
)

func TestAccResourceSpaceInstance_noNameOrGenerateName(t *testing.T) {
	project := "default"

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
				Config:      testAccResourceSpaceInstanceNoName(configPath, project),
				ExpectError: regexp.MustCompile("\"metadata.0.generate_name\": one of `metadata.0.generate_name,metadata.0.name`"),
			},
		},
	})
}

func TestAccResourceSpaceInstance_noNamespace(t *testing.T) {
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
				Config:      testAccResourceSpaceInstanceNoNamespace(configPath, name),
				ExpectError: regexp.MustCompile(`The argument "namespace" is required, but no definition was found.`),
			},
		},
	})
}
func TestAccResourceSpaceInstance_withGivenUser(t *testing.T) {
	name := names.SimpleNameGenerator.GenerateName("mycluster-")
	//name := "my-space"
	user := "admin"
	//user2 := "admin2"
	project := "default"

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
		CheckDestroy:      testAccSpaceInstanceCheckDestroy(kubeClient),
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSpaceInstanceCreateWithUser(configPath, user, project, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space_instance.test_user", "metadata.0.name", name),
					resource.TestCheckResourceAttr("loft_space_instance.test_user", "metadata.0.namespace", "loft-p-"+project),
					resource.TestCheckResourceAttr("loft_space_instance.test_user", "spec.0.owner.0.user", user),
					resource.TestCheckResourceAttr("loft_space_instance.test_user", "spec.0.owner.0.team", ""),
					checkSpaceInstance(configPath, project, name, hasUser(user)),
				),
			},
			{
				Config:            testAccResourceSpaceInstanceCreateWithUser(configPath, user, project, name),
				ResourceName:      "loft_space_instance.test_user",
				ImportState:       true,
				ImportStateVerify: true,
			},
			//{
			//	Config: testAccResourceSpaceInstanceCreateWithUser(configPath, user2, project, name),
			//	Check: resource.ComposeTestCheckFunc(
			//		resource.TestCheckResourceAttr("loft_space_instance.test_user", "metadata.0.name", name),
			//		resource.TestCheckResourceAttr("loft_space_instance.test_user", "metadata.0.namespace", "loft-p-"+project),
			//		resource.TestCheckResourceAttr("loft_space_instance.test_user", "spec.0.owner.0.user", user2),
			//		resource.TestCheckResourceAttr("loft_space_instance.test_user", "spec.0.owner.0.team", ""),
			//		checkSpaceInstance(configPath, project, name, hasUser(user2)),
			//	),
			//},
		},
	})
}

/*
	func TestAccResourceSpaceInstance_withGivenTeam(t *testing.T) {
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
						checkSpaceInstance(configPath, cluster, name, hasOwnerTeam(team)),
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
						checkSpaceInstance(configPath, cluster, name, hasOwnerTeam(team2)),
					),
				},
			},
		})
	}

	func TestAccResourceSpaceInstance_withAnnotations(t *testing.T) {
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
						checkSpaceInstance(configPath, cluster, name, hasInstanceAnnotation("some.domain/test", annotation)),
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
						checkSpaceInstance(configPath, cluster, name, hasInstanceAnnotation("some.domain/test", annotation2)),
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
						checkSpaceInstance(configPath, cluster, name, noInstanceAnnotation("some.domain/test")),
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

	func TestAccResourceSpaceInstance_withLabels(t *testing.T) {
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
						checkSpaceInstance(configPath, cluster, name, hasInstanceLabel("some.domain/test", label)),
						checkSpaceInstance(configPath, cluster, name, hasInstanceLabel(corev1.LabelMetadataName, name)),
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
						checkSpaceInstance(configPath, cluster, name, hasInstanceLabel("some.domain/test", label2)),
						checkSpaceInstance(configPath, cluster, name, hasInstanceLabel(corev1.LabelMetadataName, name)),
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
						checkSpaceInstance(configPath, cluster, name, noInstanceLabel("some.domain/test")),
						checkSpaceInstance(configPath, cluster, name, hasInstanceLabel(corev1.LabelMetadataName, name)),
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

	func TestAccResourceSpaceInstance_withInvalidSleepAfter(t *testing.T) {
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

	func TestAccResourceSpaceInstance_withSleepAfter(t *testing.T) {
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
						checkSpaceInstance(configPath, cluster, name, hasInstanceAnnotation(v1.SleepModeSleepAfterAnnotation, sleepAfterSeconds)),
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
						checkSpaceInstance(configPath, cluster, name, hasInstanceAnnotation(v1.SleepModeSleepAfterAnnotation, sleepAfter2Seconds)),
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
						checkSpaceInstance(configPath, cluster, name, noInstanceAnnotation(v1.SleepModeSleepAfterAnnotation)),
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

	func TestAccResourceSpaceInstance_withInvalidDeleteAfter(t *testing.T) {
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

	func TestAccResourceSpaceInstance_withDeleteAfter(t *testing.T) {
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
					Config: testAccResourceSpaceInstanceCreateWithDeleteAfter(configPath, cluster, name, deleteAfter),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("loft_space.test", "name", name),
						resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
						resource.TestCheckResourceAttr("loft_space.test", "user", ""),
						resource.TestCheckResourceAttr("loft_space.test", "team", ""),
						resource.TestCheckResourceAttr("loft_space.test", "delete_after", deleteAfterSeconds),
						checkSpaceInstance(configPath, cluster, name, hasInstanceAnnotation(v1.SleepModeDeleteAfterAnnotation, deleteAfterSeconds)),
					),
				},
				{
					ResourceName:      "loft_space.test",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					Config: testAccResourceSpaceInstanceCreateWithDeleteAfter(configPath, cluster, name, deleteAfter2),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("loft_space.test", "name", name),
						resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
						resource.TestCheckResourceAttr("loft_space.test", "user", ""),
						resource.TestCheckResourceAttr("loft_space.test", "team", ""),
						resource.TestCheckResourceAttr("loft_space.test", "delete_after", deleteAfter2Seconds),
						checkSpaceInstance(configPath, cluster, name, hasInstanceAnnotation(v1.SleepModeDeleteAfterAnnotation, deleteAfter2Seconds)),
					),
				},
				{
					ResourceName:      "loft_space.test",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					Config: testAccResourceSpaceInstanceCreateWithoutUserOrTeam(configPath, cluster, name),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("loft_space.test", "name", name),
						resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
						resource.TestCheckResourceAttr("loft_space.test", "user", ""),
						resource.TestCheckResourceAttr("loft_space.test", "team", ""),
						resource.TestCheckResourceAttr("loft_space.test", "delete_after", ""),
						checkSpaceInstance(configPath, cluster, name, noInstanceAnnotation(v1.SleepModeDeleteAfterAnnotation)),
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

	func TestAccResourceSpaceInstance_withSleepSchedule(t *testing.T) {
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
					Config: testAccResourceSpaceInstanceCreateWithScheduledSleep(configPath, cluster, name, sleepSchedule),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("loft_space.test", "name", name),
						resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
						resource.TestCheckResourceAttr("loft_space.test", "user", ""),
						resource.TestCheckResourceAttr("loft_space.test", "team", ""),
						resource.TestCheckResourceAttr("loft_space.test", "sleep_schedule", sleepSchedule),
						checkSpaceInstance(configPath, cluster, name, hasInstanceAnnotation(v1.SleepModeSleepScheduleAnnotation, sleepSchedule)),
					),
				},
				{
					ResourceName:      "loft_space.test",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					Config: testAccResourceSpaceInstanceCreateWithScheduledSleep(configPath, cluster, name, sleepSchedule2),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("loft_space.test", "name", name),
						resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
						resource.TestCheckResourceAttr("loft_space.test", "user", ""),
						resource.TestCheckResourceAttr("loft_space.test", "team", ""),
						resource.TestCheckResourceAttr("loft_space.test", "sleep_schedule", sleepSchedule2),
						checkSpaceInstance(configPath, cluster, name, hasInstanceAnnotation(v1.SleepModeSleepScheduleAnnotation, sleepSchedule2)),
					),
				},
				{
					ResourceName:      "loft_space.test",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					Config: testAccResourceSpaceInstanceCreateWithoutUserOrTeam(configPath, cluster, name),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("loft_space.test", "name", name),
						resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
						resource.TestCheckResourceAttr("loft_space.test", "user", ""),
						resource.TestCheckResourceAttr("loft_space.test", "team", ""),
						resource.TestCheckResourceAttr("loft_space.test", "sleep_schedule", ""),
						checkSpaceInstance(configPath, cluster, name, noInstanceAnnotation(v1.SleepModeSleepScheduleAnnotation)),
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

	func TestAccResourceSpaceInstance_withWakeupSchedule(t *testing.T) {
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
					Config: testAccResourceSpaceInstanceCreateWithScheduledWakeup(configPath, cluster, name, wakeSchedule),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("loft_space.test", "name", name),
						resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
						resource.TestCheckResourceAttr("loft_space.test", "user", ""),
						resource.TestCheckResourceAttr("loft_space.test", "team", ""),
						resource.TestCheckResourceAttr("loft_space.test", "wakeup_schedule", wakeSchedule),
						checkSpaceInstance(configPath, cluster, name, hasInstanceAnnotation(v1.SleepModeWakeupScheduleAnnotation, wakeSchedule)),
					),
				},
				{
					ResourceName:      "loft_space.test",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					Config: testAccResourceSpaceInstanceCreateWithScheduledWakeup(configPath, cluster, name, wakeSchedule2),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("loft_space.test", "name", name),
						resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
						resource.TestCheckResourceAttr("loft_space.test", "user", ""),
						resource.TestCheckResourceAttr("loft_space.test", "team", ""),
						resource.TestCheckResourceAttr("loft_space.test", "wakeup_schedule", wakeSchedule2),
						checkSpaceInstance(configPath, cluster, name, hasInstanceAnnotation(v1.SleepModeWakeupScheduleAnnotation, wakeSchedule2)),
					),
				},
				{
					ResourceName:      "loft_space.test",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					Config: testAccResourceSpaceInstanceCreateWithoutUserOrTeam(configPath, cluster, name),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("loft_space.test", "name", name),
						resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
						resource.TestCheckResourceAttr("loft_space.test", "user", ""),
						resource.TestCheckResourceAttr("loft_space.test", "team", ""),
						resource.TestCheckResourceAttr("loft_space.test", "wakeup_schedule", ""),
						checkSpaceInstance(configPath, cluster, name, noInstanceAnnotation(v1.SleepModeWakeupScheduleAnnotation)),
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

	func TestAccResourceSpaceInstance_withSpaceConstraints(t *testing.T) {
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
					Config: testAccResourceSpaceInstanceCreateWithSpaceConstraints(configPath, cluster, name, spaceConstraints),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("loft_space.test", "name", name),
						resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
						resource.TestCheckResourceAttr("loft_space.test", "user", ""),
						resource.TestCheckResourceAttr("loft_space.test", "team", ""),
						resource.TestCheckResourceAttr("loft_space.test", "space_constraints", spaceConstraints),
						checkSpaceInstance(configPath, cluster, name, hasInstanceLabel(provider.SpaceLabelSpaceConstraints, spaceConstraints)),
					),
				},
				{
					ResourceName:            "loft_space.test",
					ImportState:             true,
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"space_constraints"},
				},
				{
					Config: testAccResourceSpaceInstanceCreateWithSpaceConstraints(configPath, cluster, name, spaceConstraints2),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("loft_space.test", "name", name),
						resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
						resource.TestCheckResourceAttr("loft_space.test", "user", ""),
						resource.TestCheckResourceAttr("loft_space.test", "team", ""),
						resource.TestCheckResourceAttr("loft_space.test", "space_constraints", spaceConstraints2),
						checkSpaceInstance(configPath, cluster, name, hasInstanceLabel(provider.SpaceLabelSpaceConstraints, spaceConstraints2)),
					),
				},
				{
					ResourceName:            "loft_space.test",
					ImportState:             true,
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"space_constraints"},
				},
				{
					Config: testAccResourceSpaceInstanceCreateWithoutUserOrTeam(configPath, cluster, name),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("loft_space.test", "name", name),
						resource.TestCheckResourceAttr("loft_space.test", "cluster", cluster),
						resource.TestCheckResourceAttr("loft_space.test", "user", ""),
						resource.TestCheckResourceAttr("loft_space.test", "team", ""),
						resource.TestCheckResourceAttr("loft_space.test", "space_constraints", ""),
						checkSpaceInstance(configPath, cluster, name, noInstanceLabel(provider.SpaceLabelSpaceConstraints)),
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

	func TestAccResourceSpaceInstance_withSpaceObjects(t *testing.T) {
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
					Config: testAccResourceSpaceInstanceCreateWithSpaceObjects(configPath, cluster, name, objects1),
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
					Config: testAccResourceSpaceInstanceCreateWithSpaceObjects(configPath, cluster, name, objects2),
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
					Config: testAccResourceSpaceInstanceCreateWithSpaceObjects(configPath, cluster, name, ""),
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

	func TestAccResourceSpaceInstance_withGenerateName(t *testing.T) {
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
					Config: testAccResourceSpaceInstanceCreateWithGenerateName(configPath, cluster, prefix),
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

	func TestAccResourceSpaceInstance_withNameAndGenerateName(t *testing.T) {
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
					Config:      testAccResourceSpaceInstanceCreateWithNameAndGenerateName(configPath, cluster, name, prefix),
					ExpectError: regexp.MustCompile(`"generate_name": conflicts with name`),
				},
			},
		})
	}
*/
func testAccResourceSpaceInstanceNoName(configPath, projectName string) string {
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

	resource "loft_space_instance" "test" {
		metadata {
			namespace = "loft-p-%s"
		}
	}

`,

		configPath,
		projectName,
	)
}

func testAccResourceSpaceInstanceNoNamespace(configPath, spaceName string) string {
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

	resource "loft_space_instance" "test" {
		metadata {
			name = "%s"
		}
	}

`,

		configPath,
		spaceName,
	)
}
func checkSpaceInstance(configPath, projectName, spaceName string, pred func(obj ctrlclient.Object) error) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		apiClient, err := clientpkg.NewClientFromPath(configPath)
		if err != nil {
			return err
		}

		projectNamespace := naming.ProjectNamespace(projectName)
		managementClient, err := apiClient.Management()
		if err != nil {
			return err
		}

		spaceInstance, err := managementClient.Loft().ManagementV1().SpaceInstances(projectNamespace).Get(context.TODO(), spaceName, metav1.GetOptions{})
		if err != nil {
			return err
		}

		return pred(spaceInstance)
	}
}

func hasOwnerUser(user string) func(space *managementv1.SpaceInstance) error {
	return func(spaceInstance *managementv1.SpaceInstance) error {
		if spaceInstance.Spec.Owner == nil {
			return fmt.Errorf(
				"%s: User was not configured",
				spaceInstance.GetName(),
			)
		}

		if spaceInstance.Spec.Owner.User != user {
			return fmt.Errorf(
				"%s: User didn't match %q, got %#v",
				spaceInstance.GetName(),
				user,
				spaceInstance.Spec.Owner.User)
		}
		return nil
	}
}

func hasOwnerTeam(team string) func(space *managementv1.SpaceInstance) error {
	return func(spaceInstance *managementv1.SpaceInstance) error {
		if spaceInstance.Spec.Owner == nil {
			return fmt.Errorf(
				"%s: Team was not configured",
				spaceInstance.GetName(),
			)
		}

		if spaceInstance.Spec.Owner.Team != team {
			return fmt.Errorf(
				"%s: Team didn't match %q, got %#v",
				spaceInstance.GetName(),
				team,
				spaceInstance.Spec.Owner.Team)
		}
		return nil
	}
}
