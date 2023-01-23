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

func TestAccResourceProject_withGivenUser(t *testing.T) {
	projectName := names.SimpleNameGenerator.GenerateName("project-")
	//name := "my-space"
	user := "admin"
	//user2 := "admin2"

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
				Config: testAccResourceProjectCreateWithUser(configPath, projectName, user),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_project.test_user", "metadata.0.name", projectName),
					resource.TestCheckResourceAttr("loft_project.test_user", "spec.0.owner.0.user", user),
					resource.TestCheckResourceAttr("loft_project.test_user", "spec.0.owner.0.team", ""),
					checkProject(configPath, projectName, hasUser(user)),
				),
			},
			{
				Config:            testAccResourceProjectCreateWithUser(configPath, projectName, user),
				ResourceName:      "loft_project.test_user",
				ImportState:       true,
				ImportStateVerify: true,
			},
			//{
			//	Config: testAccResourceProjectCreateWithUser(configPath, user2, project, name),
			//	Check: resource.ComposeTestCheckFunc(
			//		resource.TestCheckResourceAttr("loft_space_instance.test_user", "metadata.0.name", name),
			//		resource.TestCheckResourceAttr("loft_space_instance.test_user", "metadata.0.namespace", "loft-p-"+project),
			//		resource.TestCheckResourceAttr("loft_space_instance.test_user", "spec.0.owner.0.user", user2),
			//		resource.TestCheckResourceAttr("loft_space_instance.test_user", "spec.0.owner.0.team", ""),
			//		checkProject(configPath, project, name, hasUser(user2)),
			//	),
			//},
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
		metadata {
		}
	}
`,
		configPath)
}

func testAccResourceProjectCreateWithUser(configPath, project, user string) string {
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

	resource "loft_project" "test_user" {
		metadata {
			name = "%s"
		}
		spec {
			owner {
				user = "%s"
			}
		}
	}
`,
		configPath,
		project,
		user,
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
			fmt.Println(projectName)
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
