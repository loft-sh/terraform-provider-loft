package tests

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	clusterv1 "github.com/loft-sh/agentapi/v2/pkg/apis/loft/cluster/v1"
	"github.com/loft-sh/loftctl/v2/pkg/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

func testAccResourceSpaceCreateWithoutUserOrTeam(configPath string, clusterName, spaceName string) string {
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

func testAccResourceSpaceCreateWithUser(configPath string, user, clusterName, spaceName string) string {
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

resource "loft_space" "test_user" {
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

func testAccResourceSpaceCreateWithTeam(configPath, team, clusterName, spaceName string) string {
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

func testAccResourceSpaceCreateWithAnnotations(configPath, clusterName, spaceName, testAnnotation string) string {
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

func testAccResourceSpaceCreateWithLabels(configPath, clusterName, spaceName, testLabel string) string {
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

func testAccResourceSpaceCreateWithSleepAfter(configPath, clusterName, spaceName, sleepAfter string) string {
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

func testAccResourceSpaceCreateWithDeleteAfter(configPath, clusterName, spaceName, deleteAfter string) string {
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

func testAccResourceSpaceCreateWithScheduledSleep(configPath, clusterName, spaceName, sleepSchedule string) string {
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

func testAccResourceSpaceCreateWithScheduledWakeup(configPath, clusterName, spaceName, wakeSchedule string) string {
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

func testAccResourceSpaceCreateWithSpaceConstraints(configPath, clusterName, spaceName, spaceConstraints string) string {
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

func testAccResourceSpaceCreateWithSpaceObjects(configPath, clusterName, spaceName, objects string) string {
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

func testAccResourceSpaceCreateWithGenerateName(configPath, clusterName, prefix string) string {
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

func testAccResourceSpaceCreateWithNameAndGenerateName(configPath, clusterName, name, prefix string) string {
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

func checkSpace(configPath, clusterName, spaceName string, pred func(obj ctrlclient.Object) error) resource.TestCheckFunc {
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

func spaceHasUser(user string) func(obj ctrlclient.Object) error {
	return func(obj ctrlclient.Object) error {
		space, ok := obj.(*clusterv1.Space)
		if !ok {
			return fmt.Errorf("object is not a space")
		}

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

func spaceHasTeam(team string) func(obj ctrlclient.Object) error {
	return func(obj ctrlclient.Object) error {
		space, ok := obj.(*clusterv1.Space)
		if !ok {
			return fmt.Errorf("object is not a space")
		}

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

func toSecondsString(durationStr string) (string, error) {
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", int(duration.Seconds())), nil
}
