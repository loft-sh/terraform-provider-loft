package tests

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	clientpkg "github.com/loft-sh/loftctl/v2/pkg/client"
	"github.com/loft-sh/loftctl/v2/pkg/client/naming"
	"github.com/loft-sh/loftctl/v2/pkg/kube"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apiserver/pkg/storage/names"
	"regexp"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
	"testing"
	"time"
)

func TestAccResourceVirtualClusterInstance_noNameOrGenerateName(t *testing.T) {
	project := "default"

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
				Config:      testAccResourceVirtualClusterInstanceNoName(configPath, project),
				ExpectError: regexp.MustCompile("\"metadata.0.generate_name\": one of `metadata.0.generate_name,metadata.0.name`"),
			},
		},
	})
}

func TestAccResourceVirtualClusterInstance_noNamespace(t *testing.T) {
	name := names.SimpleNameGenerator.GenerateName("myvcluster-")

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
				Config:      testAccResourceVirtualClusterInstanceNoNamespace(configPath, name),
				ExpectError: regexp.MustCompile(`The argument "namespace" is required, but no definition was found.`),
			},
		},
	})
}

func TestAccResourceVirtualClusterInstance_withGivenUser(t *testing.T) {
	name := names.SimpleNameGenerator.GenerateName("mycluster-")
	//name := "my-space"
	user := "admin"
	//user2 := "admin2"
	project := "default"

	kubeClient, err := newKubeClient()
	if err != nil {
		t.Fatal(err)
	}

	_, accessKey, configPath, err := loginUser(kubeClient, user)
	if err != nil {
		t.Fatal(err)
	}
	defer logout(t, kubeClient, accessKey)

	resource.Test(t, resource.TestCase{
		CheckDestroy:      testAccVirtualClusterInstanceCheckDestroy(kubeClient),
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceVirtualClusterInstanceCreateWithUser(configPath, project, name, user),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_virtual_cluster_instance.test_user", "metadata.0.name", name),
					resource.TestCheckResourceAttr("loft_virtual_cluster_instance.test_user", "metadata.0.namespace", "loft-p-"+project),
					resource.TestCheckResourceAttr("loft_virtual_cluster_instance.test_user", "spec.0.owner.0.user", user),
					resource.TestCheckResourceAttr("loft_virtual_cluster_instance.test_user", "spec.0.owner.0.team", ""),
					checkVirtualClusterInstance(configPath, project, name, hasUser(user)),
				),
			},
			{
				Config:            testAccResourceVirtualClusterInstanceCreateWithUser(configPath, project, name, user),
				ResourceName:      "loft_virtual_cluster_instance.test_user",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"metadata.0.generation",
					"metadata.0.resource_version",
				},
			},
			//{
			//	Config: testAccResourceVirtualClusterInstanceCreateWithUser(configPath, user2, project, name),
			//	Check: resource.ComposeTestCheckFunc(
			//		resource.TestCheckResourceAttr("loft_virtual_cluster_instance.test_user", "metadata.0.name", name),
			//		resource.TestCheckResourceAttr("loft_virtual_cluster_instance.test_user", "metadata.0.namespace", "loft-p-"+project),
			//		resource.TestCheckResourceAttr("loft_virtual_cluster_instance.test_user", "spec.0.owner.0.user", user2),
			//		resource.TestCheckResourceAttr("loft_virtual_cluster_instance.test_user", "spec.0.owner.0.team", ""),
			//		checkVirtualClusterInstance(configPath, project, name, hasUser(user2)),
			//	),
			//},
		},
	})
}

func testAccResourceVirtualClusterInstanceNoName(configPath, projectName string) string {
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

resource "loft_virtual_cluster_instance" "test_user" {
	metadata {
		namespace = "loft-p-%s"
	}
}
`,
		configPath,
		projectName)
}

func testAccResourceVirtualClusterInstanceNoNamespace(configPath, name string) string {
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

resource "loft_virtual_cluster_instance" "test_user" {
	metadata {
		name = "%s"
	}
}
`,
		configPath,
		name)
}

func testAccResourceVirtualClusterInstanceCreateWithUser(configPath, projectName, name, user string) string {
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

resource "loft_virtual_cluster_instance" "test_user" {
	metadata {
		namespace = "loft-p-%s"
		name = "%s"
	}

	spec {
		owner {
			user = "%s"
		}
		template {
			metadata {}
		}
	}
}`,
		configPath,
		projectName,
		name,
		user)
}

func checkVirtualClusterInstance(configPath, projectName, name string, pred func(obj ctrlclient.Object) error) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		apiClient, err := clientpkg.NewClientFromPath(configPath)
		if err != nil {
			return err
		}

		managementClient, err := apiClient.Management()
		if err != nil {
			return err
		}

		project, err := managementClient.Loft().ManagementV1().VirtualClusterInstances(naming.ProjectNamespace(projectName)).Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return err
		}

		return pred(project)
	}
}

func testAccVirtualClusterInstanceCheckDestroy(kubeClient kube.Interface) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		var vcis []string
		for _, resourceState := range s.RootModule().Resources {
			vcis = append(vcis, resourceState.Primary.ID)
		}

		for _, vci := range vcis {
			tokens := strings.Split(vci, "/")
			vciNamespace := tokens[0]
			vciName := tokens[1]
			err := wait.PollImmediate(1*time.Second, 60*time.Second, func() (bool, error) {
				_, err := kubeClient.Loft().ManagementV1().VirtualClusterInstances(vciNamespace).Get(context.TODO(), vciName, metav1.GetOptions{})
				if errors.IsNotFound(err) {
					return true, nil
				}
				return false, err
			})
			if err != nil {
				return err
			}
		}
		return nil
	}
}
