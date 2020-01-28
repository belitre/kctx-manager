package kubeconfig

import (
	"fmt"
	"os"
	"path"
	"text/tabwriter"

	"k8s.io/client-go/tools/clientcmd"

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

func AddContext(kubeconfigArg, newKubeconfigPath, newName string) error {
	kubeconfigPath, err := getKubeconfigPath(kubeconfigArg)
	if err != nil {
		return err
	}

	currentConfig, err := clientcmd.LoadFromFile(kubeconfigPath)

	if err != nil {
		return err
	}

	toAddConfig, err := clientcmd.LoadFromFile(newKubeconfigPath)

	if err != nil {
		return err
	}

	isChangeName := len(newName) > 0

	if isChangeName && len(toAddConfig.Contexts) > 1 {
		return fmt.Errorf("Error, can't rename context to %s, more than 1 context found in kubeconfig %s", newName, newKubeconfigPath)
	}

	// rename everything in the config to add to avoid problems with configs like eks
	for k, v := range toAddConfig.Contexts {
		cluster := toAddConfig.Clusters[v.Cluster]
		user := toAddConfig.AuthInfos[v.AuthInfo]

		delete(toAddConfig.Clusters, v.Cluster)
		delete(toAddConfig.AuthInfos, v.AuthInfo)

		var updatedName string
		if isChangeName {
			updatedName = newName
		} else {
			updatedName = k
		}

		v.Cluster = updatedName
		v.AuthInfo = updatedName
		toAddConfig.Contexts[updatedName] = v
		toAddConfig.Clusters[updatedName] = cluster
		toAddConfig.AuthInfos[updatedName] = user
		delete(toAddConfig.Contexts, k)
		toAddConfig.Contexts[updatedName] = v
	}

	// now we have all the new contexts with the correct names, so let's merge maps
	for k, v := range toAddConfig.Contexts {
		// but of course if we have the same context already we have to update it properly
		if currentCtx, ok := currentConfig.Contexts[k]; ok {
			// context found! let's remove it cause we are adding a new one
			delete(currentConfig.Clusters, currentCtx.Cluster)
			delete(currentConfig.AuthInfos, currentCtx.AuthInfo)
			delete(currentConfig.Contexts, k)
		}
		currentConfig.Contexts[k] = v
		currentConfig.Clusters[v.Cluster] = toAddConfig.Clusters[v.Cluster]
		currentConfig.AuthInfos[v.AuthInfo] = toAddConfig.AuthInfos[v.AuthInfo]
	}

	if err = clientcmd.WriteToFile(*currentConfig, kubeconfigPath); err != nil {
		return err
	}

	for k, _ := range toAddConfig.Contexts {
		fmt.Println(fmt.Sprintf("Context %s added/updated", k))
	}

	return nil
}

func DeleteContext(kubeconfigArg, contextName string) error {
	kubeconfigPath, err := getKubeconfigPath(kubeconfigArg)
	if err != nil {
		return err
	}

	currentConfig, err := clientcmd.LoadFromFile(kubeconfigPath)

	if err != nil {
		return err
	}

	if _, ok := currentConfig.Contexts[contextName]; !ok {
		fmt.Println(fmt.Sprintf("Context %s not found in %s", contextName, kubeconfigPath))
		return nil
	}

	deleteContext := currentConfig.Contexts[contextName]

	delete(currentConfig.Clusters, deleteContext.Cluster)
	delete(currentConfig.AuthInfos, deleteContext.AuthInfo)
	delete(currentConfig.Contexts, contextName)

	if err = clientcmd.WriteToFile(*currentConfig, kubeconfigPath); err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("Context %s deleted successfully!", contextName))

	return nil
}

func RenameContext(kubeconfigArg, contextName, newName string, isForce bool) error {
	kubeconfigPath, err := getKubeconfigPath(kubeconfigArg)
	if err != nil {
		return err
	}

	currentConfig, err := clientcmd.LoadFromFile(kubeconfigPath)

	if err != nil {
		return err
	}

	if _, ok := currentConfig.Contexts[contextName]; !ok {
		fmt.Println(fmt.Sprintf("Context %s not found in %s", contextName, kubeconfigPath))
		return nil
	}

	if _, ok := currentConfig.Contexts[newName]; ok {
		if !isForce {
			return fmt.Errorf("Error, context %s already exists in %s", newName, kubeconfigPath)
		}
	}

	renameContext := currentConfig.Contexts[contextName]
	renameUser := currentConfig.AuthInfos[renameContext.AuthInfo]
	renameCluster := currentConfig.Clusters[renameContext.Cluster]

	//rename everything
	renameContext.Cluster = newName
	renameContext.AuthInfo = newName

	currentConfig.Contexts[newName] = renameContext
	currentConfig.Clusters[newName] = renameCluster
	currentConfig.AuthInfos[newName] = renameUser

	delete(currentConfig.Contexts, contextName)
	delete(currentConfig.Clusters, contextName)
	delete(currentConfig.AuthInfos, contextName)

	if currentConfig.CurrentContext == contextName {
		currentConfig.CurrentContext = newName
	}

	if err = clientcmd.WriteToFile(*currentConfig, kubeconfigPath); err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("Context %s renamed to %s successfully!", contextName, newName))

	return nil
}

func ListContexts(kubeconfigArg string) error {
	kubeconfigPath, err := getKubeconfigPath(kubeconfigArg)
	if err != nil {
		return err
	}

	currentConfig, err := clientcmd.LoadFromFile(kubeconfigPath)

	if err != nil {
		return err
	}

	listContexts := []*ContextWithEndpoint{}

	for k, v := range currentConfig.Contexts {
		ctxWithEndpoint := &ContextWithEndpoint{
			Endpoint: currentConfig.Clusters[v.Cluster].Server,
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
