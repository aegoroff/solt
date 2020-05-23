package cmd

import (
	"github.com/spf13/afero"
	"io"
	"os"

	"github.com/spf13/cobra"
)

const (
	solutionFileExt    = ".sln"
	csharpProjectExt   = ".csproj"
	cppProjectExt      = ".vcxproj"
	packagesConfigFile = "packages.config"

	pathParamName = "path"
	diagParamName = "diag"
)

var sourcesPath string
var showMemUsage bool

var appFileSystem afero.Fs
var appWriter io.Writer

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "solt",
	Short: "SOLution Tool that analyzes Microsoft Visual Studio solutions and projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
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
	appWriter = os.Stdout
	appFileSystem = afero.NewOsFs()
	cobra.MousetrapHelpText = ""
	rootCmd.PersistentFlags().StringVarP(&sourcesPath, pathParamName, "p", "", "REQUIRED. Path to the sources folder")
	rootCmd.PersistentFlags().BoolVarP(&showMemUsage, diagParamName, "d", false, "Show memory statistic after run")
}
