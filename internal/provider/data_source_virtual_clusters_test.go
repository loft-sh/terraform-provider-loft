package provider

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"k8s.io/apiserver/pkg/storage/names"
)

func TestAccDataSourceVirtualClusters_all(t *testing.T) {
	rxPosNum := regexp.MustCompile("^[1-9][0-9]*$")
	user := "admin"
	team := "loft-admins"
	clusterName := "loft-cluster"
	namespace := "testn1s"

	virtualCluster1Name := names.SimpleNameGenerator.GenerateName("my-virtual-cluster-1-")
	annotation := names.SimpleNameGenerator.GenerateName("annotation-")

	virtualCluster2Name := names.SimpleNameGenerator.GenerateName("my-virtual-cluster-2-")
	label := names.SimpleNameGenerator.GenerateName("label-")

	virtualCluster3Name := names.SimpleNameGenerator.GenerateName("my-virtual-cluster-3-")
	values := `storage:
  size: 5Gi
`
	virtualCluster4Name := names.SimpleNameGenerator.GenerateName("objects-")
	objects := `apiVersion: v1
kind: ConfigMap
metadata:
 name: test-config-map
data:
 foo: bar
`

	kubeClient, err := newKubeClient()
	if err != nil {
		t.Error(err)
		return
	}

	loftClient, adminAccessKey, configPath, err := loginUser(kubeClient, user)
	if err != nil {
		t.Fatal(err)
	}
	defer logout(t, kubeClient, adminAccessKey)

	teamAccessKey, clusterAccess, _, err := loginTeam(kubeClient, loftClient, clusterName, team)
	if err != nil {
		t.Fatal(err)
	}
	defer logout(t, kubeClient, teamAccessKey)

	defer deleteClusterAccess(t, loftClient, clusterName, clusterAccess.GetName())

	// Create space
	if err := createSpace(configPath, clusterName, namespace); err != nil {
		t.Fatal(err)
	}
	defer func(configPath, clusterName, spaceName string) {
		if err := deleteSpace(configPath, clusterName, spaceName); err != nil {
			t.Fatal(err)
		}
	}(configPath, clusterName, namespace)

	resource.Test(t, resource.TestCase{
		CheckDestroy:      testAccVirtualClusterCheckDestroy(kubeClient),
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVirtualClustersCreate(configPath, clusterName, namespace, virtualCluster1Name, annotation, virtualCluster2Name, label, virtualCluster3Name, values, virtualCluster4Name, objects),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_objects", "cluster", clusterName),
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_objects", "name", virtualCluster4Name),
				),
			},
			{
				Config: testAccDataSourceVirtualClustersCreate(configPath, clusterName, namespace, virtualCluster1Name, annotation, virtualCluster2Name, label, virtualCluster3Name, values, virtualCluster4Name, objects) +
					testAccDataSourceVirtualClustersAll(clusterName, namespace),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.loft_virtual_clusters.all", "virtual_clusters.#", rxPosNum),
					checkVirtualClusterByName("data.loft_virtual_clusters.all", virtualCluster1Name, "name", virtualCluster1Name),
					checkVirtualClusterByName("data.loft_virtual_clusters.all", virtualCluster1Name, "cluster", clusterName),
					checkVirtualClusterByName("data.loft_virtual_clusters.all", virtualCluster2Name, "name", virtualCluster2Name),
					checkVirtualClusterByName("data.loft_virtual_clusters.all", virtualCluster2Name, "cluster", clusterName),
					checkVirtualClusterByName("data.loft_virtual_clusters.all", virtualCluster3Name, "name", virtualCluster3Name),
					checkVirtualClusterByName("data.loft_virtual_clusters.all", virtualCluster3Name, "cluster", clusterName),
					checkVirtualClusterByName("data.loft_virtual_clusters.all", virtualCluster4Name, "name", virtualCluster4Name),
					checkVirtualClusterByName("data.loft_virtual_clusters.all", virtualCluster4Name, "cluster", clusterName),
					checkVirtualClusterByName("data.loft_virtual_clusters.all", virtualCluster4Name, "objects", objects),
				),
			},
		},
	})
}

func testAccDataSourceVirtualClustersCreate(configPath, clusterName, namespace, virtualCluster1Name, annotation, virtualCluster2Name, labels, virtualCluster3Name, values, virtualCluster4Name, objects string) string {
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
	name = "%[4]s"
	cluster = "%[2]s"
	namespace = "%[3]s"
	annotations = {
		"some.domain/test" = "%[5]s"
	}
}

resource "loft_virtual_cluster" "test_labels" {
	name = "%[6]s"
	cluster = "%[2]s"
    namespace = "%[3]s"
	labels = {
		"some.domain/test" = "%[7]s"
	}
}

resource "loft_virtual_cluster" "test_values" {
	name = "%[8]s"
	cluster = "%[2]s"
    namespace = "%[3]s"
	values = <<YAML
%[9]sYAML
}

resource "loft_virtual_cluster" "test_objects" {
	name = "%[10]s"
	cluster = "%[2]s"
    namespace = "%[3]s"
	objects = <<YAML
%[11]sYAML
}
`,
		configPath,
		clusterName,
		namespace,
		virtualCluster1Name,
		annotation,
		virtualCluster2Name,
		labels,
		virtualCluster3Name,
		values,
		virtualCluster4Name,
		objects,
	)
}

func testAccDataSourceVirtualClustersAll(clusterName string, namespace string) string {
	return fmt.Sprintf(`
data "loft_virtual_clusters" "all" {
	cluster = "%s"
	namespace = "%s"
}
`,
		clusterName,
		namespace,
	)
}

func checkVirtualClusterByName(moduleName, virtualClusterName, key, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		virtualClusterPath := ""
		virtualClusterNameMatch := regexp.MustCompile(`virtual_clusters\.\d+\.name`)
		primaryModule := s.RootModule().Resources[moduleName].Primary

		for k, v := range primaryModule.Attributes {
			if virtualClusterNameMatch.MatchString(k) && v == virtualClusterName {
				tokens := strings.Split(k, ".")
				virtualClusterPath = strings.Join([]string{tokens[0], tokens[1]}, ".")
				break
			}
		}

		if virtualClusterPath == "" {
			return fmt.Errorf("virtualCluster with name %s not found", virtualClusterName)
		}

		attrKey := strings.Join([]string{virtualClusterPath, key}, ".")
		if primaryModule.Attributes[attrKey] != value {
			return fmt.Errorf(
				"%s: Attribute '%s' didn't match %q, got %#v",
				moduleName,
				attrKey,
				value,
				primaryModule.Attributes[key])
		}

		return nil
	}
}
