package tests

import (
	"fmt"
)

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

func testAccResourceSpaceInstanceCreateWithUser(configPath string, user, clusterName, spaceName string) string {
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
