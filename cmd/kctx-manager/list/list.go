package list

import (
	"github.com/belitre/kctx-manager/pkg/kubeconfig"
	"github.com/spf13/cobra"
)

func CreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Show current contexts",
		RunE: func(cmd *cobra.Command, args []string) error {
			kubeconfigArg, err := cmd.Flags().GetString("kubeconfig")
			if err != nil {
				return err
			}
			return list(kubeconfigArg)
		},
	}

	return cmd
}

func list(kubeconfigArg string) error {
	return kubeconfig.ListContexts(kubeconfigArg)
}
