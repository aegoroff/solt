package cmd

import (
	"github.com/spf13/cobra"
)

// Version defines program version
var Version = "0.9.0"

func newVersion(c conf) *cobra.Command {
	var versionCmd = &cobra.Command{
		Use:     "ver",
		Aliases: []string{"version"},
		Short:   "Print the version number of solt",
		Long:    `All software has versions. This is solt's`,
		Run: func(cmd *cobra.Command, args []string) {
			c.prn().cprint("solt v%s\n", Version)
		},
	}

	return versionCmd
}
