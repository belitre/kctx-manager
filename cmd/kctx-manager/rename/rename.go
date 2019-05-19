package rename

import (
	"fmt"

	"github.com/belitre/kctx-manager/pkg/kubeconfig"
	"github.com/spf13/cobra"
)

var isForceFlag = false

func CreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rename current_context_name new_context_name",
		Short: "Rename current_context_name to new_context_name",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("incorrect number of arguments")
			}
			kubeconfigArg, err := cmd.Flags().GetString("kubeconfig")
			if err != nil {
				return err
			}
			return rename(kubeconfigArg, args[0], args[1], isForceFlag)
		},
	}

	cmd.Flags().BoolVarP(&isForceFlag, "force", "f", false, "Forces rename. If new_context_name already exists it will be deleted.")

	return cmd
}

func rename(kubeconfigArg, currentContextName, newContextName string, isForce bool) error {
	return kubeconfig.RenameContext(kubeconfigArg, currentContextName, newContextName, isForce)
}
