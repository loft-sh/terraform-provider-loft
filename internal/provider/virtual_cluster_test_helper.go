package provider

import "fmt"

func testAccDataSourceVirtualClusterCreateWithAnnotations(configPath, clusterName, namespace, virtualClusterName, testAnnotation string) string {
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

resource "loft_virtual_cluster" "test_annotations" {
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

func testAccDataSourceVirtualClusterCreateWithLabels(configPath, clusterName, namespace, virtualClusterName, testLabel string) string {
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

resource "loft_virtual_cluster" "test_labels" {
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

func testAccDataSourceVirtualClusterCreateWithVirtualClusterObjects(configPath, clusterName, namespace, virtualClusterName, objects string) string {
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
func testAccDataSourceVirtualClusterCreateWithVirtualClusterValues(configPath, clusterName, namespace, virtualClusterName, values string) string {
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

resource "loft_virtual_cluster" "test_values" {
	name = "%s"
	cluster = "%s"
    namespace = "%s"
	values = <<YAML
%sYAML
}
`,
		configPath,
		virtualClusterName,
		clusterName,
		namespace,
		values,
	)
}
