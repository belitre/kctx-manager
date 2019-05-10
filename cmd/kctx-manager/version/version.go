package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version should be updated each time there is a new release
var (
	Version   = "v0.0.1"
	GitCommit = ""
)

const (
	programName = "kctx-manager"
)

func CreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Prints the version number of " + programName,
		Long:  `Prints the version number of ` + programName,
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
