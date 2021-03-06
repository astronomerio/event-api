package cmd

import (
	"fmt"

	"github.com/astronomerio/event-api/pkg/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of the clickstream-api server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.Version)
		fmt.Println(version.GitCommit)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
