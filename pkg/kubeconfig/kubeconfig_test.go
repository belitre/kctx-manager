package kubeconfig

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

const (
	KUBECONFIG_PATH         = "resources/kubeconfig.yaml"
	DEFAULT_KUBECONFIG_PATH = "resources/expected/00_default.yaml"
)

func initDefaultKubeconfig(t *testing.T) {
	defaultKubeconfig, err := getClustersConfig(DEFAULT_KUBECONFIG_PATH)

	assert.NoError(t, err)

	err = saveConfig(defaultKubeconfig, KUBECONFIG_PATH)

	assert.NoError(t, err)
}

func TestKctxManager(t *testing.T) {
	initDefaultKubeconfig(t)
	assertEquals(t, DEFAULT_KUBECONFIG_PATH)
	RenameContext(KUBECONFIG_PATH, "minikube", "belitre")
	assertEquals(t, "resources/expected/01_rename.yaml")
	RenameContext(KUBECONFIG_PATH, "docker-for-desktop", "blah")
	assertEquals(t, "resources/expected/02_rename.yaml")
	RenameContext(KUBECONFIG_PATH, "belitre", "minikube")
	RenameContext(KUBECONFIG_PATH, "blah", "docker-for-desktop")
	RenameContext(KUBECONFIG_PATH, "test", "lololo")
	assertEquals(t, DEFAULT_KUBECONFIG_PATH)
	AddContext(KUBECONFIG_PATH, "resources/01_add.yaml")
	assertEquals(t, "resources/expected/03_add.yaml")
	AddContext(KUBECONFIG_PATH, "resources/02_add.yaml")
	assertEquals(t, "resources/expected/04_add.yaml")
	DeleteContext(KUBECONFIG_PATH, "bobedilla")
	assertEquals(t, "resources/expected/05_delete.yaml")
	DeleteContext(KUBECONFIG_PATH, "patata")
	DeleteContext(KUBECONFIG_PATH, "rancher")
	DeleteContext(KUBECONFIG_PATH, "coolcluster")
	assertEquals(t, DEFAULT_KUBECONFIG_PATH)
}

func assertEquals(t *testing.T, expected string) {
	expectedConfigs, err := getClustersConfig(expected)

	assert.NoError(t, err)

	kubeconfig, err := getClustersConfig("resources/kubeconfig.yaml")

	assert.NoError(t, err)

	assert.Equal(t, expectedConfigs.CurrentContext, kubeconfig.CurrentContext)

	// doing this to avoid errors when the order is different
	// in the kubeconfig and the expected result
	expect := splitClusters(expectedConfigs)

	current := splitClusters(kubeconfig)

	assert.True(t, cmp.Equal(expect, current))
}
