package provider

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"k8s.io/apiserver/pkg/storage/names"
	"regexp"
	"strings"
	"testing"
)

func TestAccDataSourceVirtualClusters_all(t *testing.T) {
	rxPosNum := regexp.MustCompile("^[1-9][0-9]*$")
	user := "admin"
	team := "loft-admins"
	clusterName := "loft-cluster"
	namespace := "testn1s"
	virtualCluster1Name := names.SimpleNameGenerator.GenerateName("my-virtual-cluster-1-")
	virtualCluster2Name := names.SimpleNameGenerator.GenerateName("my-virtual-cluster-2-")
	virtualCluster3Name := names.SimpleNameGenerator.GenerateName("my-virtual-cluster-3-")

	annotation := names.SimpleNameGenerator.GenerateName("annotation-")
	label := names.SimpleNameGenerator.GenerateName("label-")
	values := `storage:
  size: 5Gi
`
	objectsName := names.SimpleNameGenerator.GenerateName("objects-")
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
				Config: testAccDataSourceVirtualClusterCreateWithAnnotations(configPath, clusterName, namespace, virtualCluster1Name, annotation),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_annotations", "cluster", clusterName),
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_annotations", "name", virtualCluster1Name),
				),
			},
			{
				Config: testAccDataSourceVirtualClusterCreateWithLabels(configPath, clusterName, namespace, virtualCluster2Name, label),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_labels", "cluster", clusterName),
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_labels", "name", virtualCluster2Name),
				),
			},
			{
				Config: testAccDataSourceVirtualClusterCreateWithVirtualClusterValues(configPath, clusterName, namespace, virtualCluster3Name, values),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_values", "cluster", clusterName),
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_values", "name", virtualCluster3Name),
				),
			},
			{
				Config: testAccDataSourceVirtualClusterCreateWithVirtualClusterObjects(configPath, clusterName, namespace, objectsName, objects),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_objects", "cluster", clusterName),
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_objects", "name", objectsName),
				),
			},
			{
				Config: testAccDataSourceVirtualClustersAll(configPath, clusterName, namespace),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.loft_virtual_clusters.all", "virtual_clusters.#", rxPosNum),
					//checkVirtualClusterByName("data.loft_virtual_clusters.all", virtualCluster1Name, "name", virtualCluster1Name),
					//checkVirtualClusterByName("data.loft_virtual_clusters.all", virtualCluster1Name, "cluster", clusterName),
					//checkVirtualClusterByName("data.loft_virtual_clusters.all", virtualCluster2Name, "name", virtualCluster2Name),
					//checkVirtualClusterByName("data.loft_virtual_clusters.all", virtualCluster2Name, "cluster", clusterName),
					//checkVirtualClusterByName("data.loft_virtual_clusters.all", virtualCluster3Name, "name", virtualCluster3Name),
					//checkVirtualClusterByName("data.loft_virtual_clusters.all", virtualCluster3Name, "cluster", clusterName),
					checkVirtualClusterByName("data.loft_virtual_clusters.all", objectsName, "name", objectsName),
					checkVirtualClusterByName("data.loft_virtual_clusters.all", objectsName, "cluster", clusterName),
					checkVirtualClusterByName("data.loft_virtual_clusters.all", objectsName, "objects", objects),
				),
			},
		},
	})
}

func testAccDataSourceVirtualClustersAll(configPath string, clusterName string, namespace string) string {
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

data "loft_virtual_clusters" "all" {
cluster = "%s"
namespace = "%s"
}
`,
		configPath,
		clusterName,
		namespace,
	)
}

func checkVirtualClusterByName(moduleName, virtualClusterName, key, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		virtualClusterPath := ""
		virtualClusterNameMatch := regexp.MustCompile(`virtual_clusters\.\d+\.name`)
		primaryModule := s.RootModule().Resources[moduleName].Primary
		//fmt.Println("primaryModule : ", primaryModule.Attributes)
		marshal, _ := json.Marshal(primaryModule)
		fmt.Println(string(marshal))
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
