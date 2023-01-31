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

func TestAccResourceSpaceInstance_noNameOrGenerateName(t *testing.T) {
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
				Config:      testAccResourceSpaceInstanceNoName(configPath, project),
				ExpectError: regexp.MustCompile("\"metadata.0.generate_name\": one of `metadata.0.generate_name,metadata.0.name`"),
			},
		},
	})
}

func TestAccResourceSpaceInstance_noNamespace(t *testing.T) {
	name := names.SimpleNameGenerator.GenerateName("mycluster-")

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
				Config:      testAccResourceSpaceInstanceNoNamespace(configPath, name),
				ExpectError: regexp.MustCompile(`The argument "namespace" is required, but no definition was found.`),
			},
		},
	})
}

func TestAccResourceSpaceInstance_withTemplate(t *testing.T) {
	name := names.SimpleNameGenerator.GenerateName("my-space-instance-")
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
		CheckDestroy:      spaceInstanceCheckDestroy(kubeClient),
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSpaceInstanceCreateWithTemplate(configPath, user, project, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space_instance.test_user", "metadata.0.name", name),
					resource.TestCheckResourceAttr("loft_space_instance.test_user", "metadata.0.namespace", "loft-p-"+project),
					resource.TestCheckResourceAttr("loft_space_instance.test_user", "spec.0.owner.0.user", user),
					resource.TestCheckResourceAttr("loft_space_instance.test_user", "spec.0.owner.0.team", ""),
					checkSpaceInstance(configPath, project, name, hasUser(user)),
				),
			},
			{
				ResourceName:      "loft_space_instance.test_user",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"metadata.0.generation",
					"metadata.0.resource_version",
				},
			},
			{
				Config: testAccResourceSpaceInstanceCreateWithTemplate(configPath, user2, project, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space_instance.test_user", "metadata.0.name", name),
					resource.TestCheckResourceAttr("loft_space_instance.test_user", "metadata.0.namespace", "loft-p-"+project),
					resource.TestCheckResourceAttr("loft_space_instance.test_user", "spec.0.owner.0.user", user2),
					resource.TestCheckResourceAttr("loft_space_instance.test_user", "spec.0.owner.0.team", ""),
					checkSpaceInstance(configPath, project, name, hasUser(user2)),
				),
			},
			{
				Config: testAccResourceSpaceInstanceCreateWithTemplate(configPath, user2, project, name) +
					testAccDataSourceSpaceInstance(project, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "metadata.0.name", name),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "metadata.0.namespace", "loft-p-"+project),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.access.0.name", "instance-admin-access"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.access.0.users.0", "admin2"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.access.0.verbs.0", "use"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.access.0.verbs.1", "get"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.access.0.verbs.2", "update"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.access.0.verbs.3", "delete"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.access.0.verbs.4", "patch"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.access.1.name", "instance-access"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.access.1.teams.0", "loft-admins"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.access.1.users.0", user2),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.access.1.users.1", user2),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.access.1.verbs.0", "use"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.access.1.verbs.1", "get"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.cluster_ref.0.cluster", "loft-cluster"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.cluster_ref.0.namespace", "loft-default-s-my-space-instance-"+name),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.description", "Terraform Managed Space Instance"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.display_name", "Terraform Managed Space Instance"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.extra_access_rules.0.cluster_role", "loft:admins"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.extra_access_rules.0.teams.0", "loft-admins"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.extra_access_rules.0.users.0", user2),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.owner.0.user", user2),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.owner.0.team", ""),
					resource.TestMatchResourceAttr("data.loft_space_instance.test_user", "spec.0.parameters", regexp.MustCompile(".+")),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.template.0.access.0.default_cluster_role", "loft-cluster-space-admin"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.template.0.access.0.rules.0.cluster_role", "loft-cluster-space-admin"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.template.0.access.0.rules.0.users.0", user2),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.template.0.apps.0.name", "cert-issuer"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.template.0.apps.0.namespace", "default"),
					resource.TestMatchResourceAttr("data.loft_space_instance.test_user", "spec.0.template.0.apps.0.parameters", regexp.MustCompile(".+")),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.template.0.apps.0.release_name", "cert-issuer"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.template.0.apps.0.version", "0.0.1"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.template.0.charts.0.insecure_skip_tls_verify", "true"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.template.0.charts.0.name", "foo1"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.template.0.charts.0.password", "foo"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.template.0.charts.0.release_name", "foo1"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.template.0.charts.0.release_namespace", "default"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.template.0.charts.0.repo_url", "https://charts.example.com"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.template.0.charts.0.timeout", "10"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.template.0.charts.0.username", "foo1"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.template.0.charts.0.version", "0.0.1"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.template.0.charts.1.insecure_skip_tls_verify", "true"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.template.0.charts.1.name", "foo2"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.template.0.charts.1.password", "foo"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.template.0.charts.1.release_name", "foo2"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.template.0.charts.1.release_namespace", "default"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.template.0.charts.1.repo_url", "https://charts.example.com"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.template.0.charts.1.timeout", "10"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.template.0.charts.1.username", "foo2"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.template.0.charts.1.version", "0.0.1"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.template.0.metadata.0.annotations.foo", "bar"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.template.0.metadata.0.labels.foo", "bar"),
					resource.TestMatchResourceAttr("data.loft_space_instance.test_user", "spec.0.template.0.objects", regexp.MustCompile(".+")),
				),
			},
		},
	})
}

func TestAccResourceSpaceInstance_withTemplateRef(t *testing.T) {
	name := names.SimpleNameGenerator.GenerateName("my-space-instance-")
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
		CheckDestroy:      spaceInstanceCheckDestroy(kubeClient),
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSpaceInstanceCreateWithTemplateRef(configPath, user, project, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space_instance.test_user", "metadata.0.name", name),
					resource.TestCheckResourceAttr("loft_space_instance.test_user", "metadata.0.namespace", "loft-p-"+project),
					resource.TestCheckResourceAttr("loft_space_instance.test_user", "spec.0.owner.0.user", user),
					resource.TestCheckResourceAttr("loft_space_instance.test_user", "spec.0.owner.0.team", ""),
					checkSpaceInstance(configPath, project, name, hasUser(user)),
				),
			},
			{
				ResourceName:      "loft_space_instance.test_user",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"metadata.0.generation",
					"metadata.0.resource_version",
				},
			},
			{
				Config: testAccResourceSpaceInstanceCreateWithTemplateRef(configPath, user2, project, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_space_instance.test_user", "metadata.0.name", name),
					resource.TestCheckResourceAttr("loft_space_instance.test_user", "metadata.0.namespace", "loft-p-"+project),
					resource.TestCheckResourceAttr("loft_space_instance.test_user", "spec.0.owner.0.user", user2),
					resource.TestCheckResourceAttr("loft_space_instance.test_user", "spec.0.owner.0.team", ""),
					checkSpaceInstance(configPath, project, name, hasUser(user2)),
				),
			},
			{
				Config: testAccResourceSpaceInstanceCreateWithTemplateRef(configPath, user2, project, name) +
					testAccDataSourceSpaceInstance(project, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "metadata.0.name", name),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "metadata.0.namespace", "loft-p-"+project),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.access.0.name", "instance-admin-access"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.access.0.users.0", "admin2"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.access.0.verbs.0", "use"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.access.0.verbs.1", "get"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.access.0.verbs.2", "update"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.access.0.verbs.3", "delete"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.access.0.verbs.4", "patch"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.access.1.name", "instance-access"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.access.1.teams.0", "loft-admins"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.access.1.users.0", user2),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.access.1.verbs.0", "use"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.access.1.verbs.1", "get"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.cluster_ref.0.cluster", "loft-cluster"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.cluster_ref.0.namespace", "loft-default-s-my-space-instance-"+name),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.description", "Terraform Managed Space Instance"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.display_name", "Terraform Managed Space Instance"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.extra_access_rules.0.cluster_role", "loft:admins"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.extra_access_rules.0.teams.0", "loft-admins"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.extra_access_rules.0.users.0", user2),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.owner.0.user", user2),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.owner.0.team", ""),
					resource.TestMatchResourceAttr("data.loft_space_instance.test_user", "spec.0.parameters", regexp.MustCompile(".+")), resource.TestMatchResourceAttr("data.loft_space_instance.test_user", "spec.0.parameters", regexp.MustCompile(".+")),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.template_ref.0.name", "example-template"),
					resource.TestCheckResourceAttr("data.loft_space_instance.test_user", "spec.0.template_ref.0.version", "0.0.0"),
				),
			},
		},
	})
}

func testAccResourceSpaceInstanceNoName(configPath, projectName string) string {
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

	resource "loft_space_instance" "test" {
		metadata {
			namespace = "loft-p-%s"
		}
		spec {}
	}

`,

		configPath,
		projectName,
	)
}

func testAccResourceSpaceInstanceNoNamespace(configPath, spaceName string) string {
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

	resource "loft_space_instance" "test" {
		metadata {
			name = "%s"
		}
	}

`,

		configPath,
		spaceName,
	)
}

func testAccResourceSpaceInstanceCreateWithTemplate(configPath string, user, projectName, spaceName string) string {
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

resource "loft_space_instance" "test_user" {
	metadata {
		name = "%[2]s"
		namespace = "loft-p-%[3]s"
	}
	spec {
		access {
			name = "instance-admin-access"
			verbs = ["use", "get", "update", "delete", "patch"]
			users = ["%[4]s"]
		}
		access {
			name = "instance-access"
			verbs = ["use", "get"]
			teams = ["loft-admins"]
			users = ["%[4]s", "%[4]s"]
		}
		cluster_ref {
			cluster = "loft-cluster"
			namespace = "loft-default-s-my-space-instance-%[2]s"
		}
		description = "Terraform Managed Space Instance"
		display_name = "Terraform Managed Space Instance"
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
`,
		configPath,
		spaceName,
		projectName,
		user,
	)
}

func testAccResourceSpaceInstanceCreateWithTemplateRef(configPath string, user, projectName, spaceName string) string {
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

resource "loft_space_instance" "test_user" {
	metadata {
		name = "%[2]s"
		namespace = "loft-p-%[3]s"
	}
	spec {
		access {
			name = "instance-admin-access"
			verbs = ["use", "get", "update", "delete", "patch"]
			users = ["%[4]s"]
		}
		access {
			name = "instance-access"
			verbs = ["use", "get"]
			teams = ["loft-admins"]
			users = ["%[4]s"]
		}
		cluster_ref {
			cluster = "loft-cluster"
			namespace = "loft-default-s-my-space-instance-%[2]s"
		}
		description = "Terraform Managed Space Instance"
		display_name = "Terraform Managed Space Instance"
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
}
`,
		configPath,
		spaceName,
		projectName,
		user,
	)
}

func testAccDataSourceSpaceInstance(projectName, spaceName string) string {
	return fmt.Sprintf(`
data "loft_space_instance" "test_user" {
	metadata {
		namespace = "loft-p-%s"
		name = "%s"
	}
}
`,
		projectName,
		spaceName,
	)
}

func checkSpaceInstance(configPath, projectName, spaceName string, pred func(obj ctrlclient.Object) error) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		apiClient, err := clientpkg.NewClientFromPath(configPath)
		if err != nil {
			return err
		}

		projectNamespace := naming.ProjectNamespace(projectName)
		managementClient, err := apiClient.Management()
		if err != nil {
			return err
		}

		spaceInstance, err := managementClient.Loft().ManagementV1().SpaceInstances(projectNamespace).Get(context.TODO(), spaceName, metav1.GetOptions{})
		if err != nil {
			return err
		}

		return pred(spaceInstance)
	}
}

func spaceInstanceCheckDestroy(kubeClient kube.Interface) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		var spaces []string
		for _, resourceState := range s.RootModule().Resources {
			spaces = append(spaces, resourceState.Primary.ID)
		}

		for _, spacePath := range spaces {
			tokens := strings.Split(spacePath, "/")
			spaceNamespace := tokens[0]
			spaceName := tokens[1]

			err := wait.PollImmediate(1*time.Second, 60*time.Second, func() (bool, error) {
				_, err := kubeClient.Loft().ManagementV1().SpaceInstances(spaceNamespace).Get(context.TODO(), spaceName, metav1.GetOptions{})
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
