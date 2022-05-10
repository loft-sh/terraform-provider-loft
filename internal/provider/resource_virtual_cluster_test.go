package provider

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	v1 "github.com/loft-sh/agentapi/v2/pkg/apis/loft/cluster/v1"
	"github.com/loft-sh/loftctl/v2/pkg/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/storage/names"
)

func TestAccResourceVirtualCluster_noName(t *testing.T) {
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

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceVirtualClusterNoName(configPath, cluster, namespace),
				ExpectError: regexp.MustCompile(`The argument "name" is required, but no definition was found.`),
			},
		},
	})
}

func TestAccResourceVirtualCluster_noCluster(t *testing.T) {
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
				Config: testAccDataSourceVirtualClusterCreateWithAnnotations(configPath, cluster, namespace, name, annotation1),
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
				Config: testAccDataSourceVirtualClusterCreateWithAnnotations(configPath, cluster, namespace, name, annotation2),
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
				Config: testAccDataSourceVirtualClusterCreateWithLabels(configPath, cluster, namespace, name, label1),
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
				Config: testAccDataSourceVirtualClusterCreateWithLabels(configPath, cluster, namespace, name, label2),
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
				Config: testAccDataSourceVirtualClusterCreateWithVirtualClusterValues(configPath, cluster, namespace, name, values),
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
				Config: testAccDataSourceVirtualClusterCreateWithVirtualClusterObjects(configPath, cluster, namespace, name, objects1),
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
				Config: testAccDataSourceVirtualClusterCreateWithVirtualClusterObjects(configPath, cluster, namespace, name, objects2),
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
				Config: testAccDataSourceVirtualClusterCreateWithVirtualClusterObjects(configPath, cluster, namespace, name, ""),
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

		virtualCluster, err := clusterClient.Agent().ClusterV1().VirtualClusters(namespace).Get(context.TODO(), virtualClusterName, metav1.GetOptions{})
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
