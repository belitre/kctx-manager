package add

import (
	"fmt"
	"os"

	"github.com/belitre/kctx-manager/pkg/kubeconfig"

	"github.com/spf13/cobra"
)

func CreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add kubeconfig_file",
		Short: "Adds the contexts defined in kubeconfig_file",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("incorrect number of arguments")
			}
			kubeconfigArg, err := cmd.Flags().GetString("kubeconfig")
			if err != nil {
				return err
			}
			return add(kubeconfigArg, args[0])
		},
	}
	return cmd
}

func add(kubeconfigArg, newKubeconfigPath string) error {
	if _, err := os.Stat(newKubeconfigPath); os.IsNotExist(err) {
		return fmt.Errorf("file %s not found", newKubeconfigPath)
	}

	err := kubeconfig.AddContext(kubeconfigArg, newKubeconfigPath)

	return err
}
