package delete

import (
	"fmt"

	"github.com/belitre/kctx-manager/pkg/kubeconfig"
	"github.com/spf13/cobra"
)

func CreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete context_name",
		Short: "Deletes the context specified in context_name",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("incorrect number of arguments")
			}
			kubeconfigFile, err := cmd.Flags().GetString("file")
			if err != nil {
				return err
			}
			return delete(kubeconfigFile, args[0])
		},
	}
	return cmd
}

func delete(kubeconfigPath, contextName string) error {
	err := kubeconfig.DeleteContext(kubeconfigPath, contextName)
	return err
}
