package main

import (
	"os"

	"github.com/belitre/kctx-manager/cmd/kctx-manager/add"
	"github.com/belitre/kctx-manager/cmd/kctx-manager/delete"
	"github.com/belitre/kctx-manager/cmd/kctx-manager/list"
	"github.com/belitre/kctx-manager/cmd/kctx-manager/rename"
	"github.com/belitre/kctx-manager/cmd/kctx-manager/version"
	"github.com/spf13/cobra"
)

var kubeconfig string

// nolint gomnd
func main() {
	rootCmd := &cobra.Command{
		Use:   "kctx-manager",
		Short: "A CLI tool to manage your $HOME/.kube/config",
	}

	rootCmd.AddCommand(version.CreateCommand())
	rootCmd.AddCommand(list.CreateCommand())
	rootCmd.AddCommand(add.CreateCommand())
	rootCmd.AddCommand(delete.CreateCommand())
	rootCmd.AddCommand(rename.CreateCommand())

	rootCmd.PersistentFlags().StringVarP(&kubeconfig, "kubeconfig", "k", "", "Path to kubeconfig file, if not provided it will default to "+
		"the value of $KUBECONFIG, and if the environment variable is not set it will default to $HOME/.kube/config")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
