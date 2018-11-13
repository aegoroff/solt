package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

const solutionFileExt = ".sln"
const csharpProjectExt = ".csproj"
const cppProjectExt = ".vcxproj"
const packagesConfigFile = "packages.config"

const pathParamName = "path"

var sourcesPath string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "solt",
	Short: "SOLution Tool that analyzes Microsoft Visual Studio solutions and projects",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.MousetrapHelpText = ""
	rootCmd.PersistentFlags().StringVarP(&sourcesPath, pathParamName, "p", "", "REQUIRED. Path to the sources folder")
}
