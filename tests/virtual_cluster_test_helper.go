package tests

import (
	"context"
	"fmt"
	agentv1 "github.com/loft-sh/agentapi/v2/pkg/apis/loft/cluster/v1"
	"github.com/loft-sh/loftctl/v2/pkg/client"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createSpace(configPath, clusterName, spaceName string) error {
	apiClient, err := client.NewClientFromPath(configPath)
	if err != nil {
		return err
	}

	clusterClient, err := apiClient.Cluster(clusterName)
	if err != nil {
		return err
	}

	space := &agentv1.Space{
		ObjectMeta: metav1.ObjectMeta{Name: spaceName},
		Spec:       agentv1.SpaceSpec{},
	}

	_, err = clusterClient.Agent().ClusterV1().Spaces().Create(context.TODO(), space, metav1.CreateOptions{})
	if err != nil && !errors.IsAlreadyExists(err) {
		return err
	}

	return nil
}

func deleteSpace(configPath, clusterName, spaceName string) error {
	apiClient, err := client.NewClientFromPath(configPath)
	if err != nil {
		return err
	}

	clusterClient, err := apiClient.Cluster(clusterName)
	if err != nil {
		return err
	}

	if err := clusterClient.Agent().ClusterV1().Spaces().Delete(context.TODO(), spaceName, metav1.DeleteOptions{}); err != nil {
		return err
	}

	return nil
}

func testAccResourceVirtualClusterCreateWithGenerateName(configPath, clusterName, namespace, generateName string) string {
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

resource "loft_virtual_cluster" "test_generate_name" {
	generate_name = "%s"
	cluster = "%s"
    namespace = "%s"
}
`,
		configPath,
		generateName,
		clusterName,
		namespace,
	)
}

func testAccResourceVirtualClusterCreateWithNameAndGenerateName(configPath, clusterName, namespace, generateName, name string) string {
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

resource "loft_virtual_cluster" "test_generate_name" {
	name = "%s"
	generate_name = "%s"
	cluster = "%s"
    namespace = "%s"
}
`,
		configPath,
		name,
		generateName,
		clusterName,
		namespace,
	)
}

func testAccResourceVirtualClusterCreateWithAnnotations(configPath, clusterName, namespace, virtualClusterName, testAnnotation string) string {
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

func testAccResourceVirtualClusterCreateWithLabels(configPath, clusterName, namespace, virtualClusterName, testLabel string) string {
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

func testAccResourceVirtualClusterCreateWithVirtualClusterObjects(configPath, clusterName, namespace, virtualClusterName, objects string) string {
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
func testAccResourceVirtualClusterCreateWithVirtualClusterValues(configPath, clusterName, namespace, virtualClusterName, values string) string {
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
