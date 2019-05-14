package rename

import (
	"fmt"

	"github.com/belitre/kctx-manager/pkg/kubeconfig"
	"github.com/spf13/cobra"
)

func CreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rename current_context_name new_context_name",
		Short: "Renames current_context_name to new_context_name",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("incorrect number of arguments")
			}
			kubeconfigArg, err := cmd.Flags().GetString("kubeconfig")
			if err != nil {
				return err
			}
			return rename(kubeconfigArg, args[0], args[1])
		},
	}
	return cmd
}

func rename(kubeconfigArg, currentContextName, newContextName string) error {
	return kubeconfig.RenameContext(kubeconfigArg, currentContextName, newContextName)
}
