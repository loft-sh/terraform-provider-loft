package tests

import (
	"context"
	"fmt"
	"regexp"
	"strings"
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

func TestAccResourceVirtualClusterTemplate_noNameOrGenerateName(t *testing.T) {
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
				Config:      testAccResourceVirtualClusterTemplateNoName(configPath),
				ExpectError: regexp.MustCompile("\"metadata.0.generate_name\": one of `metadata.0.generate_name,metadata.0.name`"),
			},
		},
	})
}

func TestAccResourceVirtualClusterTemplate_minimal(t *testing.T) {
	name := names.SimpleNameGenerator.GenerateName("my-template-")
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
		CheckDestroy:      virtualClusterTemplateCheckDestroy(kubeClient),
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceVirtualClusterTemplateMinimal(configPath, name, user),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_virtual_cluster_template.test_user", "metadata.0.name", name),
					resource.TestCheckResourceAttr("loft_virtual_cluster_template.test_user", "spec.0.owner.0.user", user),
					resource.TestCheckResourceAttr("loft_virtual_cluster_template.test_user", "spec.0.owner.0.team", ""),
					checkVirtualClusterTemplate(configPath, name, hasUser(user)),
				),
			},
			{
				ResourceName:      "loft_virtual_cluster_template.test_user",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"metadata.0.generation",
					"metadata.0.resource_version",
				},
			},
			{
				Config: testAccResourceVirtualClusterTemplateMinimal(configPath, name, user2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_virtual_cluster_template.test_user", "metadata.0.name", name),
					resource.TestCheckResourceAttr("loft_virtual_cluster_template.test_user", "spec.0.owner.0.user", user2),
					resource.TestCheckResourceAttr("loft_virtual_cluster_template.test_user", "spec.0.owner.0.team", ""),
					checkVirtualClusterTemplate(configPath, name, hasUser(user2)),
				),
			},
		},
	})
}

func TestAccResourceVirtualClusterTemplate_allProperties(t *testing.T) {
	name := names.SimpleNameGenerator.GenerateName("my-template-")
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
		CheckDestroy:      virtualClusterTemplateCheckDestroy(kubeClient),
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceVirtualClusterTemplateAllProperties(configPath, name, user),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_virtual_cluster_template.test_user", "metadata.0.name", name),
					resource.TestCheckResourceAttr("loft_virtual_cluster_template.test_user", "spec.0.owner.0.user", user),
					resource.TestCheckResourceAttr("loft_virtual_cluster_template.test_user", "spec.0.owner.0.team", ""),
					checkVirtualClusterTemplate(configPath, name, hasUser(user)),
				),
			},
			{
				ResourceName:      "loft_virtual_cluster_template.test_user",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"metadata.0.generation",
					"metadata.0.resource_version",
				},
			},
			{
				Config: testAccResourceVirtualClusterTemplateAllProperties(configPath, name, user2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("loft_virtual_cluster_template.test_user", "metadata.0.name", name),
					resource.TestCheckResourceAttr("loft_virtual_cluster_template.test_user", "spec.0.owner.0.user", user2),
					resource.TestCheckResourceAttr("loft_virtual_cluster_template.test_user", "spec.0.owner.0.team", ""),
					checkVirtualClusterTemplate(configPath, name, hasUser(user2)),
				),
			}, {
				Config: testAccResourceVirtualClusterTemplateAllProperties(configPath, name, user2) +
					testAccDataSourceVirtualClusterTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "metadata.0.name", name),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.access.0.name", "instance-admin-access"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.access.0.subresources.0", "logs"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.access.0.subresources.1", "kubeconfig"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.access.0.users.0", "admin2"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.access.0.verbs.0", "create"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.access.0.verbs.1", "use"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.access.0.verbs.2", "get"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.access.0.verbs.3", "update"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.access.0.verbs.4", "delete"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.access.0.verbs.5", "patch"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.access.1.name", "instance-access"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.access.1.subresources.0", "logs"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.access.1.subresources.1", "kubeconfig"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.access.1.teams.0", "loft-admins"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.access.1.users.0", user2),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.access.1.users.1", user2),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.access.1.verbs.0", "create"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.access.1.verbs.1", "use"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.access.1.verbs.2", "get"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.description", "Terraform Managed Virtual Cluster Instance"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.display_name", "Terraform Managed Virtual Cluster Instance"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.owner.0.user", user2),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.owner.0.team", ""),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.access.0.default_cluster_role", "loft-cluster-space-admin"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.access.0.rules.0.cluster_role", "loft-cluster-space-admin"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.access.0.rules.0.users.0", user2),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.apps.0.name", "cert-issuer"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.apps.0.namespace", "default"),
					resource.TestMatchResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.apps.0.parameters", regexp.MustCompile(".+")),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.apps.0.release_name", "cert-issuer"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.apps.0.version", "0.0.1"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.charts.0.insecure_skip_tls_verify", "true"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.charts.0.name", "foo1"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.charts.0.password", "foo"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.charts.0.release_name", "foo1"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.charts.0.release_namespace", "default"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.charts.0.repo_url", "https://charts.example.com"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.charts.0.timeout", "10"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.charts.0.username", "foo1"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.charts.0.version", "0.0.1"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.charts.1.insecure_skip_tls_verify", "true"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.charts.1.name", "foo2"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.charts.1.password", "foo"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.charts.1.release_name", "foo2"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.charts.1.release_namespace", "default"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.charts.1.repo_url", "https://charts.example.com"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.charts.1.timeout", "10"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.charts.1.username", "foo2"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.charts.1.version", "0.0.1"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.forward_token", "true"),
					resource.TestMatchResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.helm_release.0.values", regexp.MustCompile(".+")),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.helm_release.0.chart.0.name", "vcluster"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.helm_release.0.chart.0.repo", "https://charts.loft.sh"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.helm_release.0.chart.0.version", "0.14.0-beta.0"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.instance_template.0.metadata.0.annotations.fizz", "buzz"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.instance_template.0.metadata.0.labels.fuzz", "bizz"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.metadata.0.annotations.foo", "bar"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.metadata.0.labels.foo", "bar"),
					resource.TestMatchResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.objects", regexp.MustCompile(".+")),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.pro.0.enabled", "true"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.space_template.0.apps.0.name", "cert-issuer"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.space_template.0.apps.0.namespace", "default"),
					resource.TestMatchResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.space_template.0.apps.0.parameters", regexp.MustCompile(".+")),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.space_template.0.apps.0.release_name", "cert-issuer"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.space_template.0.apps.0.version", "0.0.1"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.space_template.0.charts.0.insecure_skip_tls_verify", "true"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.space_template.0.charts.0.name", "foo1"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.space_template.0.charts.0.password", "foo"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.space_template.0.charts.0.release_name", "foo1"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.space_template.0.charts.0.release_namespace", "default"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.space_template.0.charts.0.repo_url", "https://charts.example.com"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.space_template.0.charts.0.timeout", "10"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.space_template.0.charts.0.username", "foo1"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.space_template.0.charts.0.version", "0.0.1"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.space_template.0.charts.1.insecure_skip_tls_verify", "true"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.space_template.0.charts.1.name", "foo2"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.space_template.0.charts.1.password", "foo"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.space_template.0.charts.1.release_name", "foo2"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.space_template.0.charts.1.release_namespace", "default"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.space_template.0.charts.1.repo_url", "https://charts.example.com"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.space_template.0.charts.1.timeout", "10"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.space_template.0.charts.1.username", "foo2"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.space_template.0.charts.1.version", "0.0.1"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.space_template.0.metadata.0.annotations.foo", "bar"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.space_template.0.metadata.0.labels.foo", "bar"),
					resource.TestMatchResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.space_template.0.objects", regexp.MustCompile(".+")),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.workload_virtual_cluster_template.0.metadata.0.annotations.boo", "far"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.workload_virtual_cluster_template.0.metadata.0.labels.foo", "bar"),
					resource.TestMatchResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.workload_virtual_cluster_template.0.helm_release.0.values", regexp.MustCompile(".+")),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.workload_virtual_cluster_template.0.helm_release.0.chart.0.name", "vcluster"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.workload_virtual_cluster_template.0.helm_release.0.chart.0.repo", "https://charts.loft.sh"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.workload_virtual_cluster_template.0.helm_release.0.chart.0.version", "0.14.0-beta.0"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.workload_virtual_cluster_template.0.space_template.0.apps.0.name", "cert-issuer"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.workload_virtual_cluster_template.0.space_template.0.apps.0.namespace", "default"),
					resource.TestMatchResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.workload_virtual_cluster_template.0.space_template.0.apps.0.parameters", regexp.MustCompile(".+")),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.workload_virtual_cluster_template.0.space_template.0.apps.0.release_name", "cert-issuer"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.workload_virtual_cluster_template.0.space_template.0.apps.0.version", "0.0.1"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.workload_virtual_cluster_template.0.space_template.0.charts.0.insecure_skip_tls_verify", "true"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.workload_virtual_cluster_template.0.space_template.0.charts.0.name", "foo1"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.workload_virtual_cluster_template.0.space_template.0.charts.0.password", "foo"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.workload_virtual_cluster_template.0.space_template.0.charts.0.release_name", "foo1"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.workload_virtual_cluster_template.0.space_template.0.charts.0.release_namespace", "default"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.workload_virtual_cluster_template.0.space_template.0.charts.0.repo_url", "https://charts.example.com"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.workload_virtual_cluster_template.0.space_template.0.charts.0.timeout", "10"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.workload_virtual_cluster_template.0.space_template.0.charts.0.username", "foo1"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.workload_virtual_cluster_template.0.space_template.0.charts.0.version", "0.0.1"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.workload_virtual_cluster_template.0.space_template.0.charts.1.insecure_skip_tls_verify", "true"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.workload_virtual_cluster_template.0.space_template.0.charts.1.name", "foo2"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.workload_virtual_cluster_template.0.space_template.0.charts.1.password", "foo"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.workload_virtual_cluster_template.0.space_template.0.charts.1.release_name", "foo2"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.workload_virtual_cluster_template.0.space_template.0.charts.1.release_namespace", "default"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.workload_virtual_cluster_template.0.space_template.0.charts.1.repo_url", "https://charts.example.com"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.workload_virtual_cluster_template.0.space_template.0.charts.1.timeout", "10"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.workload_virtual_cluster_template.0.space_template.0.charts.1.username", "foo2"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.workload_virtual_cluster_template.0.space_template.0.charts.1.version", "0.0.1"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.workload_virtual_cluster_template.0.space_template.0.metadata.0.annotations.foo", "bar"),
					resource.TestCheckResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.workload_virtual_cluster_template.0.space_template.0.metadata.0.labels.foo", "bar"),
					resource.TestMatchResourceAttr("data.loft_virtual_cluster_template.test_user", "spec.0.template.0.workload_virtual_cluster_template.0.space_template.0.objects", regexp.MustCompile(".+")),
				),
			},
		},
	})
}

func testAccResourceVirtualClusterTemplateNoName(configPath string) string {
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

resource "loft_virtual_cluster_template" "test_user" {
	metadata {}
	spec {}
}
`,
		configPath)
}

func testAccResourceVirtualClusterTemplateMinimal(configPath, name, user string) string {
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

resource "loft_virtual_cluster_template" "test_user" {
	metadata {
		name = "%[2]s"
	}

	spec {
		owner {
			user = "%[3]s"
		}
		template {
			metadata {}
		}
	}
}`,
		configPath,
		name,
		user)
}

func testAccResourceVirtualClusterTemplateAllProperties(configPath, name, user string) string {
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

resource "loft_virtual_cluster_template" "test_user" {
	metadata {
		name = "%[2]s"
	}

	spec {
		access {
			name = "instance-admin-access"
		  	verbs = ["create", "use", "get", "update", "delete", "patch"]
			subresources = ["logs", "kubeconfig"]
		  	users = ["%[3]s"]
		}
		access {
			name = "instance-access"
		  	verbs = ["create", "use", "get"]
			subresources = ["logs", "kubeconfig"]
		  	users = ["%[3]s", "%[3]s"]
			teams = ["loft-admins"]
		}
		description = "Terraform Managed Virtual Cluster Instance"
		display_name = "Terraform Managed Virtual Cluster Instance"
		owner {
			user = "%[3]s"
		}
		template {
			access {
				default_cluster_role = "loft-cluster-space-admin"
				rules {
					users = ["%[3]s"]
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
			forward_token = true
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
			instance_template {
				metadata {
					annotations = {
						"fizz" = "buzz"
					}
					labels = {
						"fuzz" = "bizz"
					}
				}
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
			pro {
				enabled = true
			}
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
			workload_virtual_cluster_template {
				metadata {
					annotations = {
						"boo" = "far"
					}
					labels = {
						"foo" = "bar"
					}
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
	}
}`,
		configPath,
		name,
		user)
}

func testAccDataSourceVirtualClusterTemplate(instanceName string) string {
	return fmt.Sprintf(`
data "loft_virtual_cluster_template" "test_user" {
	metadata {
		name = "%s"
	}
}
`,
		instanceName,
	)
}

func checkVirtualClusterTemplate(configPath, name string, pred func(obj ctrlclient.Object) error) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		apiClient, err := clientpkg.NewClientFromPath(configPath)
		if err != nil {
			return err
		}

		managementClient, err := apiClient.Management()
		if err != nil {
			return err
		}

		project, err := managementClient.Loft().ManagementV1().VirtualClusterTemplates().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return err
		}

		return pred(project)
	}
}

func virtualClusterTemplateCheckDestroy(kubeClient kube.Interface) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		var vcis []string
		for _, resourceState := range s.RootModule().Resources {
			vcis = append(vcis, resourceState.Primary.ID)
		}

		for _, vci := range vcis {
			tokens := strings.Split(vci, "/")
			vciName := tokens[0]
			err := wait.PollImmediate(1*time.Second, 60*time.Second, func() (bool, error) {
				_, err := kubeClient.Loft().ManagementV1().VirtualClusterTemplates().Get(context.TODO(), vciName, metav1.GetOptions{})
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
