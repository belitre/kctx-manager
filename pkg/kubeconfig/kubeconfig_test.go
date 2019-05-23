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

	var err error

	err = RenameContext(KUBECONFIG_PATH, "minikube", "belitre", false)
	assert.NoError(t, err)
	assertEquals(t, "resources/expected/01_rename.yaml")
	err = RenameContext(KUBECONFIG_PATH, "docker-for-desktop", "blah", false)
	assert.NoError(t, err)
	assertEquals(t, "resources/expected/02_rename.yaml")
	err = RenameContext(KUBECONFIG_PATH, "belitre", "minikube", false)
	assert.NoError(t, err)
	err = RenameContext(KUBECONFIG_PATH, "blah", "docker-for-desktop", false)
	assert.NoError(t, err)
	err = RenameContext(KUBECONFIG_PATH, "test", "lololo", false)
	assert.NoError(t, err)
	assertEquals(t, DEFAULT_KUBECONFIG_PATH)
	err = AddContext(KUBECONFIG_PATH, "resources/01_add.yaml")
	assert.NoError(t, err)
	assertEquals(t, "resources/expected/03_add.yaml")
	err = AddContext(KUBECONFIG_PATH, "resources/02_add.yaml")
	assert.NoError(t, err)
	assertEquals(t, "resources/expected/04_add.yaml")
	err = DeleteContext(KUBECONFIG_PATH, "bobedilla")
	assert.NoError(t, err)
	assertEquals(t, "resources/expected/05_delete.yaml")
	err = DeleteContext(KUBECONFIG_PATH, "patata")
	assert.NoError(t, err)
	err = DeleteContext(KUBECONFIG_PATH, "rancher")
	assert.NoError(t, err)
	err = RenameContext(KUBECONFIG_PATH, "docker-for-desktop", "coolcluster", false)
	assert.Error(t, err)
	assertEquals(t, "resources/expected/06_rename_fail.yaml")
	err = RenameContext(KUBECONFIG_PATH, "docker-for-desktop", "coolcluster", true)
	assert.NoError(t, err)
	err = RenameContext(KUBECONFIG_PATH, "coolcluster", "docker-for-desktop", false)
	assert.NoError(t, err)
	assertEquals(t, DEFAULT_KUBECONFIG_PATH)

	initDefaultKubeconfig(t)
	assertEquals(t, DEFAULT_KUBECONFIG_PATH)
}

func assertEquals(t *testing.T, expected string) {
	expectedConfigs, err := getClustersConfig(expected)

	assert.NoError(t, err)

	kubeconfig, err := getClustersConfig(KUBECONFIG_PATH)

	assert.NoError(t, err)

	assert.Equal(t, expectedConfigs.CurrentContext, kubeconfig.CurrentContext)

	// doing this to avoid errors when the order is different
	// in the kubeconfig and the expected result
	expect := splitClusters(expectedConfigs)

	current := splitClusters(kubeconfig)

	assert.True(t, cmp.Equal(expect, current))
}
