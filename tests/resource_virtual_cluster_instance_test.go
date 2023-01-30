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
	name := names.SimpleNameGenerator.GenerateName("my-vcluster-")

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

func TestAccResourceVirtualClusterInstance_withMinimalTemplate(t *testing.T) {
	name := names.SimpleNameGenerator.GenerateName("my-vcluster-")
	user := "admin"
	user2 := "admin2"
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
		CheckDestroy:      virtualClusterInstanceCheckDestroy(kubeClient),
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceVirtualClusterInstanceMinimalWithTemplate(configPath, project, name, user),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_virtual_cluster_instance.test_user", "metadata.0.name", name),
					resource.TestCheckResourceAttr("loft_virtual_cluster_instance.test_user", "metadata.0.namespace", "loft-p-"+project),
					resource.TestCheckResourceAttr("loft_virtual_cluster_instance.test_user", "spec.0.owner.0.user", user),
					resource.TestCheckResourceAttr("loft_virtual_cluster_instance.test_user", "spec.0.owner.0.team", ""),
					checkVirtualClusterInstance(configPath, project, name, hasUser(user)),
				),
			},
			{
				ResourceName:      "loft_virtual_cluster_instance.test_user",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"metadata.0.generation",
					"metadata.0.resource_version",
				},
			},
			{
				Config: testAccResourceVirtualClusterInstanceMinimalWithTemplate(configPath, project, name, user2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_virtual_cluster_instance.test_user", "metadata.0.name", name),
					resource.TestCheckResourceAttr("loft_virtual_cluster_instance.test_user", "metadata.0.namespace", "loft-p-"+project),
					resource.TestCheckResourceAttr("loft_virtual_cluster_instance.test_user", "spec.0.owner.0.user", user2),
					resource.TestCheckResourceAttr("loft_virtual_cluster_instance.test_user", "spec.0.owner.0.team", ""),
					checkVirtualClusterInstance(configPath, project, name, hasUser(user2)),
				),
			},
		},
	})
}

func TestAccResourceVirtualClusterInstance_allWithTemplate(t *testing.T) {
	name := names.SimpleNameGenerator.GenerateName("my-vcluster-")
	user := "admin"
	user2 := "admin2"
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
		CheckDestroy:      virtualClusterInstanceCheckDestroy(kubeClient),
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceVirtualClusterInstanceAllWithTemplate(configPath, project, name, user),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_virtual_cluster_instance.test_user", "metadata.0.name", name),
					resource.TestCheckResourceAttr("loft_virtual_cluster_instance.test_user", "metadata.0.namespace", "loft-p-"+project),
					resource.TestCheckResourceAttr("loft_virtual_cluster_instance.test_user", "spec.0.owner.0.user", user),
					resource.TestCheckResourceAttr("loft_virtual_cluster_instance.test_user", "spec.0.owner.0.team", ""),
					checkVirtualClusterInstance(configPath, project, name, hasUser(user)),
				),
			},
			{
				ResourceName:      "loft_virtual_cluster_instance.test_user",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"metadata.0.generation",
					"metadata.0.resource_version",
				},
			},
			{
				Config: testAccResourceVirtualClusterInstanceAllWithTemplate(configPath, project, name, user2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_virtual_cluster_instance.test_user", "metadata.0.name", name),
					resource.TestCheckResourceAttr("loft_virtual_cluster_instance.test_user", "metadata.0.namespace", "loft-p-"+project),
					resource.TestCheckResourceAttr("loft_virtual_cluster_instance.test_user", "spec.0.owner.0.user", user2),
					resource.TestCheckResourceAttr("loft_virtual_cluster_instance.test_user", "spec.0.owner.0.team", ""),
					checkVirtualClusterInstance(configPath, project, name, hasUser(user2)),
				),
			},
		},
	})
}

func TestAccResourceVirtualClusterInstance_allWithTemplateRef(t *testing.T) {
	name := names.SimpleNameGenerator.GenerateName("my-vcluster-")
	user := "admin"
	user2 := "admin2"
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
		CheckDestroy:      virtualClusterInstanceCheckDestroy(kubeClient),
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceVirtualClusterInstanceAllWithTemplateRef(configPath, project, name, user),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_virtual_cluster_instance.test_user", "metadata.0.name", name),
					resource.TestCheckResourceAttr("loft_virtual_cluster_instance.test_user", "metadata.0.namespace", "loft-p-"+project),
					resource.TestCheckResourceAttr("loft_virtual_cluster_instance.test_user", "spec.0.owner.0.user", user),
					resource.TestCheckResourceAttr("loft_virtual_cluster_instance.test_user", "spec.0.owner.0.team", ""),
					checkVirtualClusterInstance(configPath, project, name, hasUser(user)),
				),
			},
			{
				ResourceName:      "loft_virtual_cluster_instance.test_user",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"metadata.0.generation",
					"metadata.0.resource_version",
				},
			},
			{
				Config: testAccResourceVirtualClusterInstanceAllWithTemplateRef(configPath, project, name, user2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_virtual_cluster_instance.test_user", "metadata.0.name", name),
					resource.TestCheckResourceAttr("loft_virtual_cluster_instance.test_user", "metadata.0.namespace", "loft-p-"+project),
					resource.TestCheckResourceAttr("loft_virtual_cluster_instance.test_user", "spec.0.owner.0.user", user2),
					resource.TestCheckResourceAttr("loft_virtual_cluster_instance.test_user", "spec.0.owner.0.team", ""),
					checkVirtualClusterInstance(configPath, project, name, hasUser(user2)),
				),
			},
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
	spec {}
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

func testAccResourceVirtualClusterInstanceMinimalWithTemplate(configPath, projectName, name, user string) string {
	return fmt.Sprintf(`
terraform {
	required_providers {
		loft = {
			source = "registry.terraform.io/loft-sh/loft"
		}
	}
}

provider "loft" {
	config_path = "%[1]s"
}

resource "loft_virtual_cluster_instance" "test_user" {
	metadata {
		namespace = "loft-p-%[2]s"
		name = "%[3]s"
	}

	spec {
		owner {
			user = "%[4]s"
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

func testAccResourceVirtualClusterInstanceAllWithTemplate(configPath, projectName, name, user string) string {
	return fmt.Sprintf(`
terraform {
	required_providers {
		loft = {
			source = "registry.terraform.io/loft-sh/loft"
		}
	}
}

provider "loft" {
	config_path = "%[1]s"
}

resource "loft_virtual_cluster_instance" "test_user" {
	metadata {
		namespace = "loft-p-%[2]s"
		name = "%[3]s"
	}

	spec {
		access {
			name = "instance-admin-access"
		  	verbs = ["create", "use", "get", "update", "delete", "patch"]
			subresources = ["logs", "kubeconfig"]
		  	users = ["%[4]s"]
		}
		access {
			name = "instance-access"
		  	verbs = ["create", "use", "get"]
			subresources = ["logs", "kubeconfig"]
		  	users = ["%[4]s", "%[4]s"]
			teams = ["loft-admins"]
		}
		cluster_ref {
			cluster = "loft-cluster"
    		namespace = "loft-default-v-my-vcluster"
    		virtual_cluster = "my-vcluster"
		}
		description = "Terraform Managed Virtual Cluster Instance"
		display_name = "Terraform Managed Virtual Cluster Instance"
		extra_access_rules {
			cluster_role = "loft:admins"
			users = ["%[4]s"]
			teams = ["loft-admins"]
		}
		owner {
			user = "%[4]s"
		}
		parameters = <<PARAMS
- variable: mylabelvalue
  label: vClusterStatefulSetLabelValue
  description: Please select the value for the vCluster statefulset "my-label" key
  options:
    - one
    - two
  section: Labels
PARAMS
		template {
			access {
				default_cluster_role = "loft-cluster-space-admin"
				rules {
					users = ["%[4]s"]
					cluster_role = "loft-cluster-space-admin"
				}
			}
			access_point {
				ingress {
					enabled = false
				}
			}
			apps {
				name = "cert-issuer"
				namespace = "default"
				release_name = "cert-issuer"
				version = "0.0.1"
				parameters = <<PARAMS
certIssuer:
  email: "test@test.com"
PARAMS
			}
			charts {
				insecure_skip_tls_verify = true
				name = "foo1"
				password = "foo"
				release_name = "foo1"
				release_namespace = "default"
				repo_url = "https://charts.example.com"
				timeout = 10
				username = "foo1"
				version = "0.0.1"
				wait = false
			}
			charts {
				insecure_skip_tls_verify = true
				name = "foo2"
				password = "foo"
				release_name = "foo2"
				release_namespace = "default"
				repo_url = "https://charts.example.com"
				timeout = 10
				username = "foo2"
				version = "0.0.1"
				wait = false
			}
			helm_release {
				chart {
					name = "vcluster"
					repo = "https://charts.loft.sh"
					version = "0.14.0-beta.0"
				}
				values = <<VALUES
ingress:
  enabled: false
VALUES
			}
			metadata {
				annotations = {
					"foo" = "bar"
				}
				labels = {
					"foo" = "bar"
				}
			}
			objects = <<OBJECTS
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config-map
data:
  foo: bar
  hello: world
OBJECTS
			space_template {
				apps {
					name = "cert-issuer"
					namespace = "default"
					release_name = "cert-issuer"
					version = "0.0.1"
					parameters = <<PARAMS
	certIssuer:
	  email: "test@test.com"
	PARAMS
				}
				charts {
					insecure_skip_tls_verify = true
					name = "foo1"
					password = "foo"
					release_name = "foo1"
					release_namespace = "default"
					repo_url = "https://charts.example.com"
					timeout = 10
					username = "foo1"
					version = "0.0.1"
					wait = false
				}
				charts {
					insecure_skip_tls_verify = true
					name = "foo2"
					password = "foo"
					release_name = "foo2"
					release_namespace = "default"
					repo_url = "https://charts.example.com"
					timeout = 10
					username = "foo2"
					version = "0.0.1"
					wait = false
				}
				metadata {
					annotations = {
						"foo" = "bar"
					}
					labels = {
						"foo" = "bar"
					}
				}
				objects = <<OBJECTS
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config-map
data:
  foo: bar
  hello: world
OBJECTS
			}
		}
	}
}`,
		configPath,
		projectName,
		name,
		user)
}

func testAccResourceVirtualClusterInstanceAllWithTemplateRef(configPath, projectName, name, user string) string {
	return fmt.Sprintf(`
terraform {
	required_providers {
		loft = {
			source = "registry.terraform.io/loft-sh/loft"
		}
	}
}

provider "loft" {
	config_path = "%[1]s"
}

resource "loft_virtual_cluster_instance" "test_user" {
	metadata {
		namespace = "loft-p-%[2]s"
		name = "%[3]s"
	}

	spec {
		access {
			name = "instance-admin-access"
		  	verbs = ["create", "use", "get", "update", "delete", "patch"]
			subresources = ["logs", "kubeconfig"]
		  	users = ["%[4]s"]
		}
		access {
			name = "instance-access"
		  	verbs = ["create", "use", "get"]
			subresources = ["logs", "kubeconfig"]
		  	users = ["%[4]s", "%[4]s"]
			teams = ["loft-admins"]
		}
		cluster_ref {
			cluster = "loft-cluster"
    		namespace = "loft-default-v-my-vcluster"
    		virtual_cluster = "my-vcluster"
		}
		description = "Terraform Managed Virtual Cluster Instance"
		display_name = "Terraform Managed Virtual Cluster Instance"
		extra_access_rules {
			cluster_role = "loft:admins"
			users = ["%[4]s"]
			teams = ["loft-admins"]
		}
		owner {
			user = "%[4]s"
		}
		parameters = <<PARAMS
- variable: mylabelvalue
  label: vClusterStatefulSetLabelValue
  description: Please select the value for the vCluster statefulset "my-label" key
  options:
    - one
    - two
  section: Labels
PARAMS
		template_ref {
			name = "example-template"
			version = "0.0.0"
			sync_once = false
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

func virtualClusterInstanceCheckDestroy(kubeClient kube.Interface) func(s *terraform.State) error {
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
