package delete

import (
	"fmt"

	"github.com/belitre/kctx-manager/pkg/kubeconfig"
	"github.com/spf13/cobra"
)

func CreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete context_name",
		Short: "Delete the context specified in context_name",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("incorrect number of arguments")
			}
			kubeconfigArg, err := cmd.Flags().GetString("kubeconfig")
			if err != nil {
				return err
			}
			return delete(kubeconfigArg, args[0])
		},
	}

	return cmd
}

func delete(kubeconfigArg, contextName string) error {
	return kubeconfig.DeleteContext(kubeconfigArg, contextName)
}
