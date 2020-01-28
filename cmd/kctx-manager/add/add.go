package add

import (
	"fmt"
	"os"

	"github.com/belitre/kctx-manager/pkg/kubeconfig"

	"github.com/spf13/cobra"
)

const maxArgs = 1

var nameArg string

func CreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add kubeconfig_file",
		Short: "Add the contexts defined in kubeconfig_file",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != maxArgs {
				return fmt.Errorf("incorrect number of arguments")
			}
			kubeconfigArg, err := cmd.Flags().GetString("kubeconfig")
			if err != nil {
				return err
			}
			return add(kubeconfigArg, args[0], nameArg)
		},
	}

	cmd.Flags().StringVarP(&nameArg, "name", "n", "", "name of the cluster."+
		" Use this argument if you want to rename the cluster while adding it.")

	return cmd
}

func add(kubeconfigArg, newKubeconfigPath, newName string) error {
	if _, err := os.Stat(newKubeconfigPath); os.IsNotExist(err) {
		return fmt.Errorf("file %s not found", newKubeconfigPath)
	}

	return kubeconfig.AddContext(kubeconfigArg, newKubeconfigPath, newName)
}
