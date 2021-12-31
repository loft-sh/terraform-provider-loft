package provider

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	agentv1 "github.com/loft-sh/agentapi/v2/pkg/apis/loft/cluster/v1"
	agentstoragev1 "github.com/loft-sh/agentapi/v2/pkg/apis/loft/storage/v1"
	v1 "github.com/loft-sh/api/v2/pkg/apis/management/v1"
	storagev1 "github.com/loft-sh/api/v2/pkg/apis/storage/v1"
	"github.com/loft-sh/api/v2/pkg/client/clientset_generated/clientset/scheme"
	"github.com/loft-sh/loftctl/v2/pkg/client"
	"github.com/loft-sh/loftctl/v2/pkg/kube"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var providerFactories = map[string]func() (*schema.Provider, error){
	"loft": func() (*schema.Provider, error) {
		return New("dev")(), nil
	},
}

func TestProvider(t *testing.T) {
	if err := New("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
}

func loginUser(c kube.Interface, user string) (client.Client, *v1.OwnedAccessKey, string, error) {
	uuid := uuid.NewUUID()

	accessKey, err := createUserAccessKey(c, user, string(uuid))
	if err != nil {
		return nil, nil, "", err
	}

	client, configPath, err := loginAndSaveConfigFile(accessKey.Spec.Key)
	if err != nil {
		return nil, nil, "", err
	}

	return client, accessKey, configPath, nil
}

func loginTeam(c kube.Interface, loftClient client.Client, clusterName, team string) (*v1.OwnedAccessKey, *agentv1.LocalClusterAccess, string, error) {
	teamAccess := fmt.Sprintf("%s-access", team)

	clusterAccess, err := createTeamClusterAccess(loftClient, clusterName, teamAccess, team)
	if err != nil {
		return nil, nil, "", err
	}

	uuid := uuid.NewUUID()
	accessKey, err := createTeamAccessKey(c, team, string(uuid))
	if err != nil {
		return nil, nil, "", err
	}

	_, configPath, err := loginAndSaveConfigFile(accessKey.Spec.Key)
	if err != nil {
		return nil, nil, "", err
	}

	return accessKey, clusterAccess, configPath, nil
}

func logout(c kube.Interface, accessKey *v1.OwnedAccessKey) error {
	err := deleteAccessKey(c, accessKey)
	if err != nil {
		return err
	}

	return nil
}

func newKubeClient() (kube.Interface, error) {
	kubeConfig := os.Getenv("KUBE_CONFIG")
	if kubeConfig == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		kubeConfig = filepath.Join(homeDir, ".kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		return nil, err
	}

	return kube.NewForConfig(config)
}

func createUserAccessKey(c kube.Interface, user string, accessKey string) (*v1.OwnedAccessKey, error) {
	owner, err := c.Loft().ManagementV1().Users().Get(context.TODO(), user, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	ownerAccessKeyName := owner.Spec.Username + "-terraform"

	ownedAccessKey := &v1.OwnedAccessKey{
		Spec: v1.OwnedAccessKeySpec{
			AccessKeySpec: storagev1.AccessKeySpec{
				DisplayName: "terraform-provider-loft-tests",
				User:        user,
			},
		},
	}
	ownedAccessKey.SetGenerateName(ownerAccessKeyName)
	controllerutil.SetControllerReference(owner, ownedAccessKey, scheme.Scheme)

	ownerAccessKey, err := c.Loft().ManagementV1().OwnedAccessKeys().Create(context.TODO(), ownedAccessKey, metav1.CreateOptions{})
	if err != nil && errors.IsAlreadyExists(err) {
		err := c.Loft().ManagementV1().OwnedAccessKeys().Delete(context.TODO(), ownerAccessKeyName, metav1.DeleteOptions{})
		if err != nil && errors.IsNotFound(err) {
			return nil, err
		}

		ownerAccessKey, err = c.Loft().ManagementV1().OwnedAccessKeys().Create(context.TODO(), ownedAccessKey, metav1.CreateOptions{})
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	return ownerAccessKey, nil
}

func createTeamAccessKey(c kube.Interface, team string, accessKey string) (*v1.OwnedAccessKey, error) {
	owner, err := c.Loft().ManagementV1().Teams().Get(context.TODO(), team, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	ownerAccessKeyName := owner.Spec.Username + "-terraform"

	ownedAccessKey := &v1.OwnedAccessKey{
		Spec: v1.OwnedAccessKeySpec{
			AccessKeySpec: storagev1.AccessKeySpec{
				DisplayName: "terraform-provider-loft-tests",
				Team:        team,
			},
		},
	}
	ownedAccessKey.SetGenerateName(ownerAccessKeyName)
	controllerutil.SetControllerReference(owner, ownedAccessKey, scheme.Scheme)

	ownerAccessKey, err := c.Loft().ManagementV1().OwnedAccessKeys().Create(context.TODO(), ownedAccessKey, metav1.CreateOptions{})
	if err != nil && errors.IsAlreadyExists(err) {
		err := c.Loft().ManagementV1().OwnedAccessKeys().Delete(context.TODO(), ownerAccessKeyName, metav1.DeleteOptions{})
		if err != nil && errors.IsNotFound(err) {
			return nil, err
		}

		ownerAccessKey, err = c.Loft().ManagementV1().OwnedAccessKeys().Create(context.TODO(), ownedAccessKey, metav1.CreateOptions{})
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	return ownerAccessKey, nil
}

func deleteAccessKey(c kube.Interface, accessKey *v1.OwnedAccessKey) error {
	err := c.Loft().ManagementV1().OwnedAccessKeys().Delete(context.TODO(), accessKey.GetName(), metav1.DeleteOptions{})
	if err != nil && errors.IsNotFound(err) {
		return err
	}

	return nil
}

func loginAndSaveConfigFile(accessKey string) (client.Client, string, error) {
	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, "", err
	}

	configPath := filepath.Join(tempDir, "config.json")

	loftClient, err := client.NewClientFromPath(configPath)
	if err != nil {
		return nil, "", err
	}

	err = loftClient.LoginWithAccessKey("https://localhost:9898", accessKey, true)
	if err != nil {
		return nil, "", err
	}

	err = loftClient.Save()
	if err != nil {
		return nil, "", err
	}

	return loftClient, configPath, nil
}

func createTeamClusterAccess(c client.Client, clusterName string, teamName string, teamAccess string) (*agentv1.LocalClusterAccess, error) {
	clusterClient, err := c.Cluster(clusterName)
	if err != nil {
		return nil, err
	}

	clusterAccess := &agentv1.LocalClusterAccess{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{},
		Spec: agentv1.LocalClusterAccessSpec{
			LocalClusterAccessSpec: agentstoragev1.LocalClusterAccessSpec{
				DisplayName: teamName,
				Teams:       []string{teamAccess},
			},
		},
	}
	clusterAccess.SetGenerateName(teamName)

	clusterAccess, err = clusterClient.Agent().ClusterV1().LocalClusterAccesses().Create(context.TODO(), clusterAccess, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return clusterAccess, nil
}

func deleteClusterAccess(c client.Client, clusterName string, teamName string) error {
	clusterClient, err := c.Cluster(clusterName)
	if err != nil {
		return err
	}

	err = clusterClient.Agent().ClusterV1().LocalClusterAccesses().Delete(context.TODO(), teamName, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	return nil
}
