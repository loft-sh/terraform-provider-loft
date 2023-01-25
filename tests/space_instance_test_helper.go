package tests

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/loft-sh/loftctl/v2/pkg/kube"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"strings"
	"time"
)

func testAccSpaceInstanceCheckDestroy(kubeClient kube.Interface) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		var spaces []string
		for _, resourceState := range s.RootModule().Resources {
			spaces = append(spaces, resourceState.Primary.ID)
		}

		for _, spacePath := range spaces {
			tokens := strings.Split(spacePath, "/")
			spaceNamespace := tokens[0]
			spaceName := tokens[1]

			err := wait.PollImmediate(1*time.Second, 60*time.Second, func() (bool, error) {
				_, err := kubeClient.Loft().ManagementV1().SpaceInstances(spaceNamespace).Get(context.TODO(), spaceName, metav1.GetOptions{})
				if errors.IsNotFound(err) {
					return true, nil
				}
				return false, err
			})
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccResourceSpaceInstanceCreateWithoutUserOrTeam(configPath string, clusterName, spaceName string) string {
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
}
`,
		configPath,
		spaceName,
		clusterName,
	)
}

func testAccResourceSpaceInstanceCreateWithUser(configPath string, user, projectName, spaceName string) string {
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

resource "loft_space_instance" "test_user" {
	metadata {
		name = "%s"
		namespace = "loft-p-%s"
	}
	spec {
		owner {
			user = "%s"
		}
		template {
			metadata {}
		}
	}
}
`,
		configPath,
		spaceName,
		projectName,
		user,
	)
}

func testAccResourceSpaceInstanceCreateWithTeam(configPath, team, clusterName, spaceName string) string {
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

resource "loft_space" "test_team" {
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

func testAccResourceSpaceInstanceCreateWithAnnotations(configPath, clusterName, spaceName, testAnnotation string) string {
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
	annotations = {
		"some.domain/test" = "%s"
	}
}
`,
		configPath,
		spaceName,
		clusterName,
		testAnnotation,
	)
}

func testAccResourceSpaceInstanceCreateWithLabels(configPath, clusterName, spaceName, testLabel string) string {
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
	labels = {
		"some.domain/test" = "%s"
	}
}
`,
		configPath,
		spaceName,
		clusterName,
		testLabel,
	)
}

func testAccResourceSpaceInstanceCreateWithSleepAfter(configPath, clusterName, spaceName, sleepAfter string) string {
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
	sleep_after = "%s"
}
`,
		configPath,
		spaceName,
		clusterName,
		sleepAfter,
	)
}

func testAccResourceSpaceInstanceCreateWithDeleteAfter(configPath, clusterName, spaceName, deleteAfter string) string {
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
	delete_after = "%s"
}
`,
		configPath,
		spaceName,
		clusterName,
		deleteAfter,
	)
}

func testAccResourceSpaceInstanceCreateWithScheduledSleep(configPath, clusterName, spaceName, sleepSchedule string) string {
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
	sleep_schedule = "%s"
}
`,
		configPath,
		spaceName,
		clusterName,
		sleepSchedule,
	)
}

func testAccResourceSpaceInstanceCreateWithScheduledWakeup(configPath, clusterName, spaceName, wakeSchedule string) string {
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
	wakeup_schedule = "%s"
}
`,
		configPath,
		spaceName,
		clusterName,
		wakeSchedule,
	)
}

func testAccResourceSpaceInstanceCreateWithSpaceConstraints(configPath, clusterName, spaceName, spaceConstraints string) string {
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
	space_constraints = "%s"
}
`,
		configPath,
		spaceName,
		clusterName,
		spaceConstraints,
	)
}

func testAccResourceSpaceInstanceCreateWithSpaceObjects(configPath, clusterName, spaceName, objects string) string {
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

resource "loft_space" "test_objects" {
	name = "%s"
	cluster = "%s"
	objects = <<YAML
%sYAML
}
`,
		configPath,
		spaceName,
		clusterName,
		objects,
	)
}

func testAccResourceSpaceInstanceCreateWithGenerateName(configPath, clusterName, prefix string) string {
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

resource "loft_space" "test_generate_name" {
	cluster = "%s"
	generate_name = "%s"
}
`,
		configPath,
		clusterName,
		prefix,
	)
}

func testAccResourceSpaceInstanceCreateWithNameAndGenerateName(configPath, clusterName, name, prefix string) string {
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

resource "loft_space" "test_name_and_generate_name" {
	cluster = "%s"
	name = "%s"
	generate_name = "%s"
}
`,
		configPath,
		clusterName,
		name,
		prefix,
	)
}
