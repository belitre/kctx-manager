package kubeconfig

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"text/tabwriter"

	v1 "k8s.io/client-go/tools/clientcmd/api/v1"

	"github.com/ghodss/yaml"
	homedir "github.com/mitchellh/go-homedir"
)

const kubeconfigEnvVar = "KUBECONFIG"
const defaultKubeconfig = "/.kube/config"

type ContextWithEndpoint struct {
	Name     string
	Endpoint string
}

func getKubeconfigPath(kubeconfigArg string) (string, error) {
	if len(kubeconfigArg) > 0 {
		return kubeconfigArg, nil
	}

	envVarValue := os.Getenv(kubeconfigEnvVar)

	if len(envVarValue) > 0 {
		return envVarValue, nil
	}

	home, err := homedir.Dir()
	if err != nil {
		return "", fmt.Errorf("Error while getting user home dir, error was: %s", err)
	}

	return path.Join(home, defaultKubeconfig), nil
}

func getContextsWithEndpoint(kubeconfigPath string) ([]*ContextWithEndpoint, error) {
	currentConfig, err := getClustersConfig(kubeconfigPath)

	if err != nil {
		return nil, err
	}

	mapClusterContext := map[string]string{}

	for _, v := range currentConfig.Contexts {
		mapClusterContext[v.Context.Cluster] = v.Name
	}

	listContexts := []*ContextWithEndpoint{}
	for _, v := range currentConfig.Clusters {
		if ctx, ok := mapClusterContext[v.Name]; ok {
			ctxWithEndpoint := &ContextWithEndpoint{
				Endpoint: v.Cluster.Server,
				Name:     ctx,
			}
			listContexts = append(listContexts, ctxWithEndpoint)
		}
	}

	return listContexts, nil
}

func getClustersConfig(kubeconfigPath string) (*v1.Config, error) {
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

func AddContext(kubeconfigArg, newKubeconfigPath string) error {
	kubeconfigPath, err := getKubeconfigPath(kubeconfigArg)
	if err != nil {
		return err
	}

	currentConfig, err := getClustersConfig(kubeconfigPath)
	if err != nil {
		return err
	}

	toAddConfig, err := getClustersConfig(newKubeconfigPath)
	if err != nil {
		return err
	}

	mapNewClusters := splitClusters(toAddConfig)

	if len(mapNewClusters) == 0 {
		fmt.Println("No contexts added/updated.")
		return nil
	}

	mapCurrentClusters := splitClusters(currentConfig)

	for n, c := range mapNewClusters {
		mapCurrentClusters[n] = c
	}

	users := []v1.NamedAuthInfo{}
	clusters := []v1.NamedCluster{}
	contexts := []v1.NamedContext{}
	for _, c := range mapCurrentClusters {
		users = append(users, c.AuthInfos...)
		clusters = append(clusters, c.Clusters...)
		contexts = append(contexts, c.Contexts...)
	}

	currentConfig.Clusters = clusters
	currentConfig.AuthInfos = users
	currentConfig.Contexts = contexts

	if err = saveConfig(currentConfig, kubeconfigPath); err != nil {
		return err
	}

	for k := range mapNewClusters {
		fmt.Println(fmt.Sprintf("Context %s added/updated", k))
	}

	return nil
}

func saveConfig(config *v1.Config, kubeconfigPath string) error {
	bytes, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("Error while Marshaling result, error was: %s", err)
	}

	if err = ioutil.WriteFile(kubeconfigPath, bytes, 0644); err != nil {
		return fmt.Errorf("Error while saving file %s, error was: %s", kubeconfigPath, err)
	}

	return nil
}

func RenameContext(kubeconfigArg, contextName, newName string, isForce bool) error {
	kubeconfigPath, err := getKubeconfigPath(kubeconfigArg)
	if err != nil {
		return err
	}

	currentConfig, err := getClustersConfig(kubeconfigPath)
	if err != nil {
		return err
	}

	mapCurrentClusters := splitClusters(currentConfig)

	if _, ok := mapCurrentClusters[contextName]; !ok {
		fmt.Println(fmt.Sprintf("Context %s not found in %s", contextName, kubeconfigPath))
		return nil
	}

	if _, ok := mapCurrentClusters[newName]; ok {
		if isForce {
			delete(mapCurrentClusters, newName)
		} else {
			return fmt.Errorf("Error, context %s already exists in %s", newName, kubeconfigPath)
		}
	}

	users := []v1.NamedAuthInfo{}
	clusters := []v1.NamedCluster{}
	contexts := []v1.NamedContext{}

	for n, c := range mapCurrentClusters {
		if n == contextName {
			c.Contexts[0].Name = newName
		}
		users = append(users, c.AuthInfos...)
		clusters = append(clusters, c.Clusters...)
		contexts = append(contexts, c.Contexts...)
	}

	currentConfig.Clusters = clusters
	currentConfig.AuthInfos = users
	currentConfig.Contexts = contexts

	if currentConfig.CurrentContext == contextName {
		currentConfig.CurrentContext = newName
	}

	if err = saveConfig(currentConfig, kubeconfigPath); err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("Context %s renamed to %s successfully!", contextName, newName))

	return nil
}

func DeleteContext(kubeconfigArg, contextName string) error {
	kubeconfigPath, err := getKubeconfigPath(kubeconfigArg)
	if err != nil {
		return err
	}

	currentConfig, err := getClustersConfig(kubeconfigPath)
	if err != nil {
		return err
	}

	mapCurrentClusters := splitClusters(currentConfig)

	if _, ok := mapCurrentClusters[contextName]; !ok {
		fmt.Println(fmt.Sprintf("Context %s not found in %s", contextName, kubeconfigPath))
		return nil
	}

	users := []v1.NamedAuthInfo{}
	clusters := []v1.NamedCluster{}
	contexts := []v1.NamedContext{}

	for n, c := range mapCurrentClusters {
		if n != contextName {
			users = append(users, c.AuthInfos...)
			clusters = append(clusters, c.Clusters...)
			contexts = append(contexts, c.Contexts...)
		}
	}

	currentConfig.Clusters = clusters
	currentConfig.AuthInfos = users
	currentConfig.Contexts = contexts

	if currentConfig.CurrentContext == contextName {
		currentConfig.CurrentContext = ""
	}

	if err = saveConfig(currentConfig, kubeconfigPath); err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("Context %s deleted successfully!", contextName))

	return nil
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

func ListContexts(kubeconfigArg string) error {
	kubeconfigPath, err := getKubeconfigPath(kubeconfigArg)
	if err != nil {
		return err
	}

	currentConfig, err := getClustersConfig(kubeconfigPath)
	if err != nil {
		return err
	}

	mapCurrentClusters := splitClusters(currentConfig)

	listContexts := []*ContextWithEndpoint{}
	for k, v := range mapCurrentClusters {
		ctxWithEndpoint := &ContextWithEndpoint{
			Endpoint: v.Clusters[0].Cluster.Server,
			Name:     k,
		}
		listContexts = append(listContexts, ctxWithEndpoint)
	}

	printContexts(listContexts, kubeconfigPath)

	return nil
}

func printContexts(contexts []*ContextWithEndpoint, kubeconfigPath string) {
	if len(contexts) == 0 {
		fmt.Println(fmt.Sprintf("No contexts found in %s", kubeconfigPath))
		return
	}

	w := new(tabwriter.Writer)

	w.Init(os.Stdout, 15, 8, 1, '\t', 0)
	defer w.Flush()

	fmt.Fprintf(w, "\n %s\t%s\t", "Context", "Endpoint")
	fmt.Fprintf(w, "\n %s\t%s\t", "-------", "--------")

	for _, v := range contexts {
		fmt.Fprintf(w, "\n %s\t%s\t", v.Name, v.Endpoint)
	}

	fmt.Fprintf(w, "\n\n")
}
