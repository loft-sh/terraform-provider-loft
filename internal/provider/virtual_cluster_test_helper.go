package provider

import "fmt"

func testAccDataSourceVirtualClusterCreate_withAnnotations(configPath, clusterName, namespace, virtualClusterName, testAnnotation string) string {
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

resource "loft_virtual_cluster" "test" {
	name = "%s"
	cluster = "%s"
    namespace = "%s"
	annotations = {
		"some.domain/test" = "%s"
	}
}
`,
		configPath,
		virtualClusterName,
		clusterName,
		namespace,
		testAnnotation,
	)
}

func testAccDataSourceVirtualClusterCreate_withLabels(configPath, clusterName, namespace, virtualClusterName, testLabel string) string {
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

resource "loft_virtual_cluster" "test" {
	name = "%s"
	cluster = "%s"
    namespace = "%s"
	labels = {
		"some.domain/test" = "%s"
	}
}
`,
		configPath,
		virtualClusterName,
		clusterName,
		namespace,
		testLabel,
	)
}

func testAccDataSourceVirtualClusterCreate_withVirtualClusterConstraints(configPath, clusterName, namespace, virtualClusterName, virtualClusterConstraints string) string {
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

resource "loft_virtual_cluster" "test" {
	name = "%s"
	cluster = "%s"
    namespace = "%s"
}
`,
		configPath,
		virtualClusterName,
		clusterName,
		namespace,
	)
}

func testAccDataSourceVirtualClusterCreate_withVirtualClusterObjects(configPath, clusterName, namespace, virtualClusterName, objects string) string {
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

resource "loft_virtual_cluster" "test_objects" {
	name = "%s"
	cluster = "%s"
    namespace = "%s"
	objects = <<YAML
%sYAML
}
`,
		configPath,
		virtualClusterName,
		clusterName,
		namespace,
		objects,
	)
}
