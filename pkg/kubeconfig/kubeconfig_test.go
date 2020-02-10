package kubeconfig

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"k8s.io/client-go/tools/clientcmd"

	clientcmdv1 "k8s.io/client-go/tools/clientcmd/api/v1"
	v1 "k8s.io/client-go/tools/clientcmd/api/v1"
)

const (
	kubeconfigPath        = "resources/kubeconfig.yaml"
	defaultKubeconfigPath = "resources/expected/00_default.yaml"
)

func initDefaultKubeconfig(t *testing.T) {
	defaultKubeconfig, err := clientcmd.LoadFromFile(defaultKubeconfigPath)

	assert.NoError(t, err)

	err = clientcmd.WriteToFile(*defaultKubeconfig, kubeconfigPath)

	assert.NoError(t, err)
}

func TestKctxManager(t *testing.T) {
	initDefaultKubeconfig(t)
	assertEquals(t, defaultKubeconfigPath)

	var err error

	err = RenameContext(kubeconfigPath, "minikube", "belitre", false)
	assert.NoError(t, err)
	assertEquals(t, "resources/expected/01_rename.yaml")

	err = RenameContext(kubeconfigPath, "docker-for-desktop", "blah", false)
	assert.NoError(t, err)
	assertEquals(t, "resources/expected/02_rename.yaml")

	err = RenameContext(kubeconfigPath, "belitre", "minikube", false)
	assert.NoError(t, err)

	err = RenameContext(kubeconfigPath, "blah", "docker-for-desktop", false)
	assert.NoError(t, err)

	err = RenameContext(kubeconfigPath, "test", "lololo", false)
	assert.NoError(t, err)
	assertEquals(t, defaultKubeconfigPath)

	err = AddContext(kubeconfigPath, "resources/01_add.yaml", "")
	assert.NoError(t, err)
	assertEquals(t, "resources/expected/03_add.yaml")

	err = AddContext(kubeconfigPath, "resources/02_add.yaml", "")
	assert.NoError(t, err)
	assertEquals(t, "resources/expected/04_add.yaml")

	err = DeleteContext(kubeconfigPath, "bobedilla")
	assert.NoError(t, err)
	assertEquals(t, "resources/expected/05_delete.yaml")

	err = DeleteContext(kubeconfigPath, "patata")
	assert.NoError(t, err)

	err = DeleteContext(kubeconfigPath, "rancher")
	assert.NoError(t, err)

	err = RenameContext(kubeconfigPath, "docker-for-desktop", "coolcluster", false)
	assert.Error(t, err)
	assertEquals(t, "resources/expected/06_rename_fail.yaml")

	err = RenameContext(kubeconfigPath, "docker-for-desktop", "coolcluster", true)
	assert.NoError(t, err)

	err = RenameContext(kubeconfigPath, "coolcluster", "docker-for-desktop", false)
	assert.NoError(t, err)
	assertEquals(t, defaultKubeconfigPath)
	initDefaultKubeconfig(t)
	assertEquals(t, defaultKubeconfigPath)
}

func getClustersConfig(kubeconfigPath string) (*clientcmdv1.Config, error) {
	fileContent, err := ioutil.ReadFile(kubeconfigPath)

	if err != nil {
		return nil, fmt.Errorf("Error while reading file %s, error was: %s", kubeconfigPath, err)
	}

	clustersConfig := &v1.Config{}

	err = yaml.Unmarshal(fileContent, clustersConfig)
	if err != nil {
		return nil, fmt.Errorf("Error while unmarshalling file %s, error was: %s", kubeconfigPath, err)
	}

	return clustersConfig, nil
}

func splitClusters(config *v1.Config) map[string]*v1.Config {
	mapConfigs := map[string]*v1.Config{}

	mapClusters := map[string]v1.NamedCluster{}
	mapUsers := map[string]v1.NamedAuthInfo{}

	for _, v := range config.Clusters {
		mapClusters[v.Name] = v
	}

	for _, v := range config.AuthInfos {
		mapUsers[v.Name] = v
	}

	for _, v := range config.Contexts {
		newConfig := &v1.Config{
			Contexts: []v1.NamedContext{
				v,
			},
			AuthInfos: []v1.NamedAuthInfo{
				mapUsers[v.Context.AuthInfo],
			},
			Clusters: []v1.NamedCluster{
				mapClusters[v.Context.Cluster],
			},
		}
		mapConfigs[v.Name] = newConfig
	}

	return mapConfigs
}

func assertEquals(t *testing.T, expected string) {
	expectedConfigs, err := getClustersConfig(expected)

	assert.NoError(t, err)

	kubeconfig, err := getClustersConfig(kubeconfigPath)

	assert.NoError(t, err)

	assert.Equal(t, expectedConfigs.CurrentContext, kubeconfig.CurrentContext)

	// doing this to avoid errors when the order is different
	// in the kubeconfig and the expected result
	expect := splitClusters(expectedConfigs)

	current := splitClusters(kubeconfig)

	assert.True(t, cmp.Equal(expect, current))
}
