package tests

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	clientpkg "github.com/loft-sh/loftctl/v2/pkg/client"
	"github.com/loft-sh/loftctl/v2/pkg/kube"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apiserver/pkg/storage/names"
	"regexp"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
	"time"
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
