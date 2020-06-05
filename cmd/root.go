package cmd

import (
	"fmt"
	"github.com/spf13/afero"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"
)

const (
	pathParamName = "path"
	diagParamName = "diag"
)

var sourcesPath string
var diag bool

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
	start := time.Now()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
	elapsed := time.Since(start)

	if diag {
		_, _ = fmt.Fprintf(appWriter, "Working time: %v\n", elapsed)
	}
}

func init() {
	appWriter = os.Stdout
	appFileSystem = afero.NewOsFs()
	cobra.MousetrapHelpText = ""
	rootCmd.PersistentFlags().StringVarP(&sourcesPath, pathParamName, "p", "", "REQUIRED. Path to the sources folder")
	rootCmd.PersistentFlags().BoolVarP(&diag, diagParamName, "d", false, "Show application diagnostic after run")
}
