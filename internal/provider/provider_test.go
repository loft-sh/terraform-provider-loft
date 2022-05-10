package provider

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	agentv1 "github.com/loft-sh/agentapi/v2/pkg/apis/loft/cluster/v1"
	agentstoragev1 "github.com/loft-sh/agentapi/v2/pkg/apis/loft/storage/v1"
	storagev1 "github.com/loft-sh/api/v2/pkg/apis/storage/v1"
	"github.com/loft-sh/api/v2/pkg/client/clientset_generated/clientset/scheme"
	"github.com/loft-sh/loftctl/v2/pkg/client"
	"github.com/loft-sh/loftctl/v2/pkg/kube"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/apimachinery/pkg/util/wait"
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

func testAccSpaceCheckDestroy(client kube.Interface) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		var spaces []string
		for _, resource := range s.RootModule().Resources {
			spaces = append(spaces, resource.Primary.ID)
		}

		for _, spacePath := range spaces {
			tokens := strings.Split(spacePath, "/")
			spaceName := tokens[1]

			err := wait.PollImmediate(1*time.Second, 60*time.Second, func() (bool, error) {
				_, err := client.Agent().ClusterV1().Spaces().Get(context.TODO(), spaceName, metav1.GetOptions{})
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

func loginUser(c kube.Interface, user string) (client.Client, *storagev1.AccessKey, string, error) {
	newUUID := uuid.NewUUID()

	accessKey, err := createUserAccessKey(c, user, string(newUUID))
	if err != nil {
		return nil, nil, "", err
	}

	loftClient, configPath, err := loginAndSaveConfigFile(accessKey.Spec.Key)
	if err != nil {
		return nil, nil, "", err
	}

	return loftClient, accessKey, configPath, nil
}

func loginTeam(c kube.Interface, loftClient client.Client, clusterName, team string) (*storagev1.AccessKey, *agentv1.LocalClusterAccess, string, error) {
	teamAccess := fmt.Sprintf("%s-access", team)

	clusterAccess, err := createTeamClusterAccess(loftClient, clusterName, teamAccess, team)
	if err != nil {
		return nil, nil, "", err
	}

	newUUID := uuid.NewUUID()
	accessKey, err := createTeamAccessKey(c, team, string(newUUID))
	if err != nil {
		return nil, nil, "", err
	}

	_, configPath, err := loginAndSaveConfigFile(accessKey.Spec.Key)
	if err != nil {
		return nil, nil, "", err
	}

	return accessKey, clusterAccess, configPath, nil
}

func logout(c kube.Interface, accessKey *storagev1.AccessKey) error {
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

func createUserAccessKey(c kube.Interface, user string, key string) (*storagev1.AccessKey, error) {
	owner, err := c.Loft().StorageV1().Users().Get(context.TODO(), user, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	accessKeyName := owner.Spec.Username + "-terraform"

	accessKey := &storagev1.AccessKey{
		Spec: storagev1.AccessKeySpec{
			DisplayName: "terraform-provider-loft-tests",
			Type:        storagev1.AccessKeyTypeLogin,
			Key:         key,
			User:        user,
		},
	}
	accessKey.SetGenerateName(accessKeyName)
	_ = controllerutil.SetControllerReference(owner, accessKey, scheme.Scheme)

	accessKey, err = c.Loft().StorageV1().AccessKeys().Create(context.TODO(), accessKey, metav1.CreateOptions{})
	if err != nil && errors.IsAlreadyExists(err) {
		err := c.Loft().StorageV1().AccessKeys().Delete(context.TODO(), accessKeyName, metav1.DeleteOptions{})
		if err != nil && errors.IsNotFound(err) {
			return nil, err
		}

		accessKey, err = c.Loft().StorageV1().AccessKeys().Create(context.TODO(), accessKey, metav1.CreateOptions{})
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	return accessKey, nil
}

func createTeamAccessKey(c kube.Interface, team string, key string) (*storagev1.AccessKey, error) {
	owner, err := c.Loft().StorageV1().Teams().Get(context.TODO(), team, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	accessKeyName := owner.Spec.Username + "-terraform"

	accessKey := &storagev1.AccessKey{
		Spec: storagev1.AccessKeySpec{
			DisplayName: "terraform-provider-loft-tests",
			Type:        storagev1.AccessKeyTypeLogin,
			Key:         key,
			Team:        team,
		},
	}
	accessKey.SetGenerateName(accessKeyName)
	_ = controllerutil.SetControllerReference(owner, accessKey, scheme.Scheme)

	accessKey, err = c.Loft().StorageV1().AccessKeys().Create(context.TODO(), accessKey, metav1.CreateOptions{})
	if err != nil && errors.IsAlreadyExists(err) {
		err := c.Loft().StorageV1().AccessKeys().Delete(context.TODO(), accessKeyName, metav1.DeleteOptions{})
		if err != nil && errors.IsNotFound(err) {
			return nil, err
		}

		accessKey, err = c.Loft().StorageV1().AccessKeys().Create(context.TODO(), accessKey, metav1.CreateOptions{})
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	return accessKey, nil
}

func deleteAccessKey(c kube.Interface, accessKey *storagev1.AccessKey) error {
	err := c.Loft().StorageV1().AccessKeys().Delete(context.TODO(), accessKey.GetName(), metav1.DeleteOptions{})
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

	err = loftClient.LoginWithAccessKey("https://localhost:8080", accessKey, true)
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
