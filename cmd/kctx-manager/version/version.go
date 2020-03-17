package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version should be updated each time there is a new release
var (
	Version   = "canary"
	GitCommit = ""
)

const (
	programName = "kctx-manager"
)

func CreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of " + programName,
		Long:  `Print the version number of ` + programName,
		Run: func(cmd *cobra.Command, args []string) {
			printVersion()
		},
	}

	return cmd
}

func printVersion() {
	fmt.Printf("%s-%s\n", programName, Version)

	if len(GitCommit) > 0 {
		fmt.Printf("git commit: %s\n", GitCommit)
	}
}
