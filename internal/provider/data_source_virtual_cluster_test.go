package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"k8s.io/apiserver/pkg/storage/names"
)

func TestAccDataSourceVirtualCluster_Annotations(t *testing.T) {
	user := "admin"
	cluster := "loft-cluster"
	virtualClusterName := names.SimpleNameGenerator.GenerateName("my-virtual-cluster-")
	namespace := "demo1"
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
		if err := deleteSpace(configPath, clusterName, namespace); err != nil {
			t.Fatal(err)
		}
	}(configPath, cluster, namespace)

	resource.Test(t, resource.TestCase{
		CheckDestroy:      testAccVirtualClusterCheckDestroy(kubeClient),
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceVirtualClusterCreateWithAnnotations(configPath, cluster, namespace, virtualClusterName, "annotations-1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_virtual_cluster.test_annotations", "cluster", cluster),
					resource.TestMatchResourceAttr("loft_virtual_cluster.test_annotations", "name", regexp.MustCompile(`^my-virtual-cluster\-.*`)),
				),
			},
			{
				Config: testAccResourceVirtualClusterCreateWithAnnotations(configPath, cluster, namespace, virtualClusterName, "annotations-1") +
					testAccDataSourceVirtualClusterRead(cluster, namespace, virtualClusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.loft_virtual_cluster.test", "cluster", cluster),
					resource.TestMatchResourceAttr("data.loft_virtual_cluster.test", "name", regexp.MustCompile(`^my-virtual-cluster\-.*`)),
				),
			},
		},
	})
}

func testAccDataSourceVirtualClusterRead(clusterName, namespace, virtualClusterName string) string {
	return fmt.Sprintf(`
data "loft_virtual_cluster" "test" {
	cluster = "%s"
	name = "%s"
	namespace = "%s"
}
`,
		clusterName,
		virtualClusterName,
		namespace,
	)
}
