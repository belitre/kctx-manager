package list

import (
	"github.com/belitre/kctx-manager/pkg/kubeconfig"
	"github.com/belitre/kctx-manager/pkg/tools"
	"github.com/spf13/cobra"
)

func CreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Shows current contexts",
		RunE: func(cmd *cobra.Command, args []string) error {
			kubeconfigPath, err := cmd.Flags().GetString("file")
			if err != nil {
				return err
			}
			return list(kubeconfigPath)
		},
	}
	return cmd
}

func list(kubeconfigPath string) error {
	contexts, err := kubeconfig.GetContextsWithEndpoint(kubeconfigPath)
	if err != nil {
		return err
	}

	tools.PrintContexts(contexts)

	return nil
}
