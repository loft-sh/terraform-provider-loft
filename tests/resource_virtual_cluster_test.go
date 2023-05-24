package tests

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	v1 "github.com/loft-sh/agentapi/v3/pkg/apis/loft/storage/v1"
	"github.com/loft-sh/loftctl/v3/pkg/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/storage/names"
)

func TestAccResourceVirtualCluster_noNameOrGenerateName(t *testing.T) {
	cluster := "loft-cluster"
	namespace := names.SimpleNameGenerator.GenerateName("namespace-")

	kubeClient, err := newKubeClient()
	if err != nil {
		t.Fatal(err)
	}

	_, accessKey, configPath, err := loginUser(kubeClient, "admin")
	if err != nil {
		t.Fatal(err)
	}

	defer logout(t, kubeClient, accessKey)

	// Create space
	if err := createSpace(configPath, cluster, namespace); err != nil {
		t.Fatal(err)
	}
	defer func(configPath, clusterName, spaceName string) {
		if err := deleteSpace(configPath, clusterName, spaceName); err != nil {
			t.Fatal(err)
		}
	}(configPath, cluster, namespace)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceVirtualClusterNoName(configPath, cluster, namespace),
				ExpectError: regexp.MustCompile(`Required value: name or generateName is required`),
			},
		},
	})
}

func TestAccResourceVirtualCluster_noProject(t *testing.T) {
	name := names.SimpleNameGenerator.GenerateName("name-")
	namespace := names.SimpleNameGenerator.GenerateName("namespace-")

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
				Config:      testAccResourceVirtualClusterNoCluster(configPath, namespace, name),
				ExpectError: regexp.MustCompile(`The argument "cluster" is required, but no definition was found.`),
			},
		},
	})
}

func TestAccResourceVirtualCluster_noNamespace(t *testing.T) {
	name := names.SimpleNameGenerator.GenerateName("name-")
	cluster := names.SimpleNameGenerator.GenerateName("mycluster-")

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
				Config:      testAccResourceVirtualClusterNoNamespace(configPath, cluster, name),
				ExpectError: regexp.MustCompile(`The argument "namespace" is required, but no definition was found.`),
			},
		},
	})
}

func TestAccResourceVirtualCluster_withGenerateName(t *testing.T) {
	generateName := "name-"
	namespace := names.SimpleNameGenerator.GenerateName("namespace-")

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

	// Create space
	if err := createSpace(configPath, cluster, namespace); err != nil {
		t.Fatal(err)
	}
	defer func(configPath, clusterName, spaceName string) {
		if err := deleteSpace(configPath, clusterName, spaceName); err != nil {
			t.Fatal(err)
		}
	}(configPath, cluster, namespace)

	resource.Test(t, resource.TestCase{
		CheckDestroy:      testAccVirtualClusterCheckDestroy(kubeClient),
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceVirtualClusterCreateWithGenerateName(configPath, cluster, namespace, generateName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_generate_name", "generate_name", generateName),
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_generate_name", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_generate_name", "namespace", namespace),
				),
			},
			{
				ResourceName:      "loft_virtual_cluster.test_generate_name",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceVirtualCluster_withNameAndGenerateName(t *testing.T) {
	generateName := "myspace-"
	name := names.SimpleNameGenerator.GenerateName("myspace-")
	namespace := names.SimpleNameGenerator.GenerateName("namespace-")

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

	// Create space
	if err := createSpace(configPath, cluster, namespace); err != nil {
		t.Fatal(err)
	}
	defer func(configPath, clusterName, spaceName string) {
		if err := deleteSpace(configPath, clusterName, spaceName); err != nil {
			t.Fatal(err)
		}
	}(configPath, cluster, namespace)

	resource.Test(t, resource.TestCase{
		CheckDestroy:      testAccVirtualClusterCheckDestroy(kubeClient),
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceVirtualClusterCreateWithNameAndGenerateName(configPath, cluster, namespace, generateName, name),
				ExpectError: regexp.MustCompile(`"generate_name": conflicts with name`),
			},
		},
	})
}

func TestAccResourceVirtualCluster_withAnnotations(t *testing.T) {
	name := names.SimpleNameGenerator.GenerateName("name-")
	namespace := names.SimpleNameGenerator.GenerateName("namespace-")

	annotation1 := names.SimpleNameGenerator.GenerateName("annotation-")
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

	// Create space
	if err := createSpace(configPath, cluster, namespace); err != nil {
		t.Fatal(err)
	}
	defer func(configPath, clusterName, spaceName string) {
		if err := deleteSpace(configPath, clusterName, spaceName); err != nil {
			t.Fatal(err)
		}
	}(configPath, cluster, namespace)

	resource.Test(t, resource.TestCase{
		CheckDestroy:      testAccVirtualClusterCheckDestroy(kubeClient),
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceVirtualClusterCreateWithAnnotations(configPath, cluster, namespace, name, annotation1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_annotations", "name", name),
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_annotations", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_annotations", "namespace", namespace),
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_annotations", "annotations.some.domain/test", annotation1),
					checkVirtualCluster(configPath, cluster, namespace, name, hasAnnotationVC("some.domain/test", annotation1)),
				),
			},
			{
				ResourceName:      "loft_virtual_cluster.test_annotations",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccResourceVirtualClusterCreateWithAnnotations(configPath, cluster, namespace, name, annotation2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_annotations", "name", name),
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_annotations", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_annotations", "namespace", namespace),
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_annotations", "annotations.some.domain/test", annotation2),
					checkVirtualCluster(configPath, cluster, namespace, name, hasAnnotationVC("some.domain/test", annotation2)),
				),
			},
			{
				ResourceName:      "loft_virtual_cluster.test_annotations",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceVirtualCluster_withLabels(t *testing.T) {
	name := names.SimpleNameGenerator.GenerateName("name-")
	namespace := names.SimpleNameGenerator.GenerateName("namespace-")
	label1 := names.SimpleNameGenerator.GenerateName("label-")
	label2 := names.SimpleNameGenerator.GenerateName("label-")
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

	// Create space
	if err := createSpace(configPath, cluster, namespace); err != nil {
		t.Fatal(err)
	}
	defer func(configPath, clusterName, spaceName string) {
		if err := deleteSpace(configPath, clusterName, spaceName); err != nil {
			t.Fatal(err)
		}
	}(configPath, cluster, namespace)

	resource.Test(t, resource.TestCase{
		CheckDestroy:      testAccVirtualClusterCheckDestroy(kubeClient),
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceVirtualClusterCreateWithLabels(configPath, cluster, namespace, name, label1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_labels", "name", name),
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_labels", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_labels", "namespace", namespace),
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_labels", "labels.some.domain/test", label1),
					checkVirtualCluster(configPath, cluster, namespace, name, hasLabelVC("some.domain/test", label1)),
				),
			},
			{
				ResourceName:      "loft_virtual_cluster.test_labels",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccResourceVirtualClusterCreateWithLabels(configPath, cluster, namespace, name, label2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_labels", "name", name),
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_labels", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_labels", "namespace", namespace),
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_labels", "labels.some.domain/test", label2),
					checkVirtualCluster(configPath, cluster, namespace, name, hasLabelVC("some.domain/test", label2)),
				),
			},
			{
				ResourceName:      "loft_virtual_cluster.test_labels",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceVirtualCluster_withValues(t *testing.T) {
	name := names.SimpleNameGenerator.GenerateName("name-")
	namespace := names.SimpleNameGenerator.GenerateName("namespace-")
	cluster := "loft-cluster"
	user := "admin"
	values := `storage:
  size: 5Gi
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

	// Create space
	if err := createSpace(configPath, cluster, namespace); err != nil {
		t.Fatal(err)
	}
	defer func(configPath, clusterName, spaceName string) {
		if err := deleteSpace(configPath, clusterName, spaceName); err != nil {
			t.Fatal(err)
		}
	}(configPath, cluster, namespace)

	resource.Test(t, resource.TestCase{
		CheckDestroy:      testAccVirtualClusterCheckDestroy(kubeClient),
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceVirtualClusterCreateWithVirtualClusterValues(configPath, cluster, namespace, name, values),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_values", "name", name),
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_values", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_values", "namespace", namespace),
				),
			},
			{
				ResourceName:      "loft_virtual_cluster.test_values",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceVirtualCluster_withVirtualClusterObjects(t *testing.T) {
	name := names.SimpleNameGenerator.GenerateName("name-")
	namespace := names.SimpleNameGenerator.GenerateName("namespace-")

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

	// Create space
	if err := createSpace(configPath, cluster, namespace); err != nil {
		t.Fatal(err)
	}
	defer func(configPath, clusterName, spaceName string) {
		if err := deleteSpace(configPath, clusterName, spaceName); err != nil {
			t.Fatal(err)
		}
	}(configPath, cluster, namespace)

	resource.Test(t, resource.TestCase{
		CheckDestroy:      testAccVirtualClusterCheckDestroy(kubeClient),
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceVirtualClusterCreateWithVirtualClusterObjects(configPath, cluster, namespace, name, objects1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_objects", "name", name),
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_objects", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_objects", "namespace", namespace),
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_objects", "objects", objects1),
				),
			},
			{
				ResourceName:      "loft_virtual_cluster.test_objects",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccResourceVirtualClusterCreateWithVirtualClusterObjects(configPath, cluster, namespace, name, objects2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_objects", "name", name),
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_objects", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_objects", "namespace", namespace),
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_objects", "objects", objects2),
				),
			},
			{
				ResourceName:      "loft_virtual_cluster.test_objects",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccResourceVirtualClusterCreateWithVirtualClusterObjects(configPath, cluster, namespace, name, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_objects", "name", name),
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_objects", "cluster", cluster),
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_objects", "namespace", namespace),
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_objects", "objects", ""),
				),
			},
			{
				ResourceName:      "loft_virtual_cluster.test_objects",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceVirtualClusterNoName(configPath, clusterName, namespace string) string {
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
	namespace = "%s"
	cluster = "%s"
}
`,
		configPath,
		namespace,
		clusterName,
	)
}

func testAccResourceVirtualClusterNoNamespace(configPath, clusterName, virtualClusterName string) string {
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
	cluster = "%s"
	name = "%s"

}
`,
		configPath,
		clusterName,
		virtualClusterName,
	)
}

func testAccResourceVirtualClusterNoCluster(configPath, namespace, virtualClusterName string) string {
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
	namespace = "%s"
}
`,
		configPath,
		namespace,
		virtualClusterName,
	)
}

func checkVirtualCluster(configPath, clusterName, namespace, virtualClusterName string, pred func(virtualCluster *v1.VirtualCluster) error) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		apiClient, err := client.NewClientFromPath(configPath)
		if err != nil {
			return err
		}

		clusterClient, err := apiClient.Cluster(clusterName)
		if err != nil {
			return err
		}

		virtualCluster, err := clusterClient.Agent().StorageV1().VirtualClusters(namespace).Get(context.TODO(), virtualClusterName, metav1.GetOptions{})
		if err != nil {
			return err
		}

		return pred(virtualCluster)
	}
}

func hasAnnotationVC(annotation, value string) func(virtualCluster *v1.VirtualCluster) error {
	return func(virtualCluster *v1.VirtualCluster) error {
		if virtualCluster.GetAnnotations()[annotation] != value {
			return fmt.Errorf(
				"%s: Annotation '%s' didn't match %q, got %#v",
				virtualCluster.GetName(),
				annotation,
				value,
				virtualCluster.GetLabels()[annotation])
		}
		return nil
	}
}

func hasLabelVC(label, value string) func(virtualCluster *v1.VirtualCluster) error {
	return func(virtualCluster *v1.VirtualCluster) error {
		if virtualCluster.GetLabels()[label] != value {
			return fmt.Errorf(
				"%s: Label '%s' didn't match %q, got %#v",
				virtualCluster.GetName(),
				label,
				value,
				virtualCluster.GetLabels()[label])
		}
		return nil
	}
}
