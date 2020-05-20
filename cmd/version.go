package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version defines program version
var Version = "0.4.0"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"ver"},
	Short:   "Print the version number of solt",
	Long:    `All software has versions. This is solt's`,
	Run: func(cmd *cobra.Command, args []string) {
		_, _ = fmt.Fprintf(appWriter, "solt v%s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
