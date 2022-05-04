package provider

import "fmt"

func testAccDataSourceSpaceCreate_withoutUserOrTeam(configPath string, clusterName, spaceName string) string {
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

func testAccDataSourceSpaceCreate_withUser(configPath string, user, clusterName, spaceName string) string {
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

func testAccDataSourceSpaceCreate_withTeam(configPath, team, clusterName, spaceName string) string {
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

func testAccDataSourceSpaceCreate_withAnnotations(configPath, clusterName, spaceName, testAnnotation string) string {
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

func testAccDataSourceSpaceCreate_withLabels(configPath, clusterName, spaceName, testLabel string) string {
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

func testAccDataSourceSpaceCreate_withSleepAfter(configPath, clusterName, spaceName string, sleepAfter int) string {
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
	sleep_after = %d
}
`,
		configPath,
		spaceName,
		clusterName,
		sleepAfter,
	)
}

func testAccDataSourceSpaceCreate_withDeleteAfter(configPath, clusterName, spaceName string, deleteAfter int) string {
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
	delete_after = %d
}
`,
		configPath,
		spaceName,
		clusterName,
		deleteAfter,
	)
}

func testAccDataSourceSpaceCreate_withScheduledSleep(configPath, clusterName, spaceName, sleepSchedule string) string {
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

func testAccDataSourceSpaceCreate_withScheduledWakeup(configPath, clusterName, spaceName, wakeSchedule string) string {
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

func testAccDataSourceSpaceCreate_withSpaceConstraints(configPath, clusterName, spaceName, spaceConstraints string) string {
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

func testAccDataSourceSpaceCreate_withSpaceObjects(configPath, clusterName, spaceName, objects string) string {
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
