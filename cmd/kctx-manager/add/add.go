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
		Short: "Adds the clusters defined in kubeconfig_file",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("incorrect number of arguments")
			}
			kubeconfigFile, err := cmd.Flags().GetString("file")
			if err != nil {
				return err
			}
			return add(kubeconfigFile, args[0])
		},
	}
	return cmd
}

func add(kubeconfigFile, newClustersPath string) error {
	if _, err := os.Stat(newClustersPath); os.IsNotExist(err) {
		return fmt.Errorf("file %s not found", newClustersPath)
	}

	err := kubeconfig.AddClusters(kubeconfigFile, newClustersPath)

	return err
}
