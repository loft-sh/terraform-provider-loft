package tests

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	clientpkg "github.com/loft-sh/loftctl/v3/pkg/client"
	"github.com/loft-sh/loftctl/v3/pkg/kube"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apiserver/pkg/storage/names"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func TestAccResourceProject_noNameOrGenerateName(t *testing.T) {
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
				Config:      testAccResourceProjectNoName(configPath),
				ExpectError: regexp.MustCompile("\"metadata.0.generate_name\": one of `metadata.0.generate_name,metadata.0.name`"),
			},
		},
	})
}

func TestAccResourceProject_allProperties(t *testing.T) {
	projectName := names.SimpleNameGenerator.GenerateName("project-")
	user := "admin"
	user2 := "admin2"

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
		CheckDestroy:      testAccProjectCheckDestroy(kubeClient),
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceProjectCreateAllProperties(configPath, projectName, user, 10),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_project.test_user", "metadata.0.name", projectName),
					resource.TestCheckResourceAttr("loft_project.test_user", "spec.0.owner.0.user", user),
					resource.TestCheckResourceAttr("loft_project.test_user", "spec.0.owner.0.team", ""),
					resource.TestCheckResourceAttr("loft_project.test_user", `spec.0.quotas.0.project.spaceinstances`, "10"),
					resource.TestCheckResourceAttr("loft_project.test_user", `spec.0.quotas.0.project.virtualclusterinstances`, "10"),
					resource.TestCheckResourceAttr("loft_project.test_user", `spec.0.quotas.0.user.spaceinstances`, "10"),
					resource.TestCheckResourceAttr("loft_project.test_user", `spec.0.quotas.0.user.virtualclusterinstances`, "10"),
					checkProject(configPath, projectName, hasUser(user)),
				),
			},
			{
				ResourceName:      "loft_project.test_user",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccResourceProjectCreateAllProperties(configPath, projectName, user2, 20),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_project.test_user", "metadata.0.name", projectName),
					resource.TestCheckResourceAttr("loft_project.test_user", "spec.0.owner.0.user", user2),
					resource.TestCheckResourceAttr("loft_project.test_user", "spec.0.owner.0.team", ""),
					resource.TestCheckResourceAttr("loft_project.test_user", `spec.0.quotas.0.project.spaceinstances`, "20"),
					resource.TestCheckResourceAttr("loft_project.test_user", `spec.0.quotas.0.project.virtualclusterinstances`, "20"),
					resource.TestCheckResourceAttr("loft_project.test_user", `spec.0.quotas.0.user.spaceinstances`, "20"),
					resource.TestCheckResourceAttr("loft_project.test_user", `spec.0.quotas.0.user.virtualclusterinstances`, "20"),
					resource.TestCheckResourceAttr("loft_project.test_user", "spec.0.owner.0.team", ""),
					checkProject(configPath, projectName, hasUser(user2)),
				),
			},
			{
				Config: testAccResourceProjectCreateAllProperties(configPath, projectName, user2, 20) +
					testAccDataSourceProjectRead(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.loft_project.test_user", "metadata.0.name", projectName),
					resource.TestCheckResourceAttr("data.loft_project.test_user", "spec.0.access.0.name", "loft-admin-access"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", "spec.0.access.0.subresources.0", "*"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", "spec.0.access.0.users.0", "admin2"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", "spec.0.access.0.verbs.0", "get"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", "spec.0.access.0.verbs.1", "update"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", "spec.0.access.0.verbs.2", "patch"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", "spec.0.access.0.verbs.3", "delete"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", "spec.0.access.1.name", "loft-access"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", "spec.0.access.1.subresources.0", "members"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", "spec.0.access.1.subresources.1", "clusters"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", "spec.0.access.1.subresources.2", "templates"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", "spec.0.access.1.subresources.3", "chartinfo"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", "spec.0.access.1.subresources.4", "charts"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", "spec.0.access.1.users.0", "*"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", "spec.0.access.1.verbs.0", "get"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", "spec.0.allowed_clusters.0.name", "*"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", "spec.0.allowed_clusters.1.name", "loft-cluster"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", "spec.0.allowed_templates.0.group", "storage.loft.sh"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", "spec.0.allowed_templates.0.kind", "VirtualClusterTemplate"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", "spec.0.allowed_templates.0.name", "*"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", "spec.0.allowed_templates.1.group", "storage.loft.sh"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", "spec.0.allowed_templates.1.kind", "SpaceTemplate"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", "spec.0.allowed_templates.1.name", "*"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", "spec.0.allowed_templates.2.group", "storage.loft.sh"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", "spec.0.allowed_templates.2.is_default", "true"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", "spec.0.allowed_templates.2.kind", "VirtualClusterTemplate"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", "spec.0.allowed_templates.2.name", "isolated-vcluster"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", "spec.0.argo_c_d.0.cluster", "loft-cluster"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", "spec.0.argo_c_d.0.enabled", "true"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", "spec.0.argo_c_d.0.namespace", "argocd"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", "spec.0.argo_c_d.0.sso.0.assigned_roles.0", "role:admin"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", "spec.0.argo_c_d.0.sso.0.enabled", "true"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", "spec.0.argo_c_d.0.sso.0.host", "https://my-argocd-domain.com"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", "spec.0.description", "Terraform Managed Project"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", "spec.0.display_name", "Terraform Managed Project"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", "spec.0.members.0.cluster_role", "loft-management-project-user"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", "spec.0.members.0.group", "storage.loft.sh"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", "spec.0.members.0.kind", "User"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", "spec.0.members.0.name", "*"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", "spec.0.namespace_pattern.0.space", "{{.Values.loft.project}}-v-{{.Values.loft.name}}"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", "spec.0.namespace_pattern.0.virtual_cluster", "{{.Values.loft.project}}-v-{{.Values.loft.name}}"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", "spec.0.owner.0.user", user2),
					resource.TestCheckResourceAttr("data.loft_project.test_user", "spec.0.owner.0.team", ""),
					resource.TestCheckResourceAttr("data.loft_project.test_user", `spec.0.quotas.0.project.spaceinstances`, "20"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", `spec.0.quotas.0.project.virtualclusterinstances`, "20"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", `spec.0.quotas.0.user.spaceinstances`, "20"),
					resource.TestCheckResourceAttr("data.loft_project.test_user", `spec.0.quotas.0.user.virtualclusterinstances`, "20"),
				),
			},
		},
	})
}

func testAccResourceProjectNoName(configPath string) string {
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

	resource "loft_project" "test" {
		metadata {}
		spec {}
	}
`,
		configPath)
}

func testAccResourceProjectCreateAllProperties(configPath, project, user string, quotaCount int) string {
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

	resource "loft_project" "test_user" {
		metadata {
			name = "%[2]s"
		}
		spec {
			access {
				name = "loft-admin-access"
				verbs = ["get", "update", "patch", "delete"]
				subresources = ["*"]
				users = ["%[3]s"]
			}
			access {
				name = "loft-access"
				verbs = ["get"]
				subresources = ["members", "clusters", "templates", "chartinfo", "charts"]
				users = ["*"]
			}
			allowed_clusters {
				name = "*"
			}
			allowed_clusters {
				name = "loft-cluster"
			}
			allowed_templates {
				kind = "VirtualClusterTemplate"
				group = "storage.loft.sh"
				name = "*"
			}
			allowed_templates {
				kind = "SpaceTemplate"
				group = "storage.loft.sh"
				name ="*"
			}
			allowed_templates {
				kind = "VirtualClusterTemplate"
				group = "storage.loft.sh"
				name = "isolated-vcluster"
				is_default = true
			}
			argo_c_d {
				enabled = true
				cluster = "loft-cluster"
				namespace = "argocd"
				sso {
					enabled = true
					host = "https://my-argocd-domain.com"
					assigned_roles = ["role:admin"]
				}
			}
			description = "Terraform Managed Project"
			display_name = "Terraform Managed Project"
			members {
				kind = "User"
				group = "storage.loft.sh"
				name = "*"
				cluster_role = "loft-management-project-user"
			}
			namespace_pattern {
				space = "{{.Values.loft.project}}-v-{{.Values.loft.name}}"
				virtual_cluster = "{{.Values.loft.project}}-v-{{.Values.loft.name}}"
			}
			owner {
				user = "%[3]s"
			}
			quotas {
				project = {
				  "spaceinstances" = "%[4]d"
				  "virtualclusterinstances" = "%[4]d"
				}
				user = {
				  "spaceinstances" = "%[4]d"
				  "virtualclusterinstances" = "%[4]d"
				}
			}
		}
	}
`,
		configPath,
		project,
		user,
		quotaCount,
	)
}

func testAccDataSourceProjectRead(project string) string {
	return fmt.Sprintf(`
data "loft_project" "test_user" {
	metadata {
		name = "%s"
	}
}
`,
		project,
	)
}

func checkProject(configPath, projectName string, pred func(obj ctrlclient.Object) error) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		apiClient, err := clientpkg.NewClientFromPath(configPath)
		if err != nil {
			return err
		}

		managementClient, err := apiClient.Management()
		if err != nil {
			return err
		}

		project, err := managementClient.Loft().ManagementV1().Projects().Get(context.TODO(), projectName, metav1.GetOptions{})
		if err != nil {
			return err
		}

		return pred(project)
	}
}

func testAccProjectCheckDestroy(kubeClient kube.Interface) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		var projects []string
		for _, resourceState := range s.RootModule().Resources {
			projects = append(projects, resourceState.Primary.ID)
		}

		for _, projectName := range projects {
			err := wait.PollImmediate(1*time.Second, 60*time.Second, func() (bool, error) {
				_, err := kubeClient.Loft().ManagementV1().Projects().Get(context.TODO(), projectName, metav1.GetOptions{})
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
