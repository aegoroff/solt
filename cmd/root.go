package cmd

import (
	"github.com/spf13/afero"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var sourcesPath string
var diag bool

var appFileSystem afero.Fs
var appPrinter printer

func newRoot() *cobra.Command {
	return &cobra.Command{
		Use:   "solt",
		Short: "SOLution Tool that analyzes Microsoft Visual Studio solutions and projects",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(args ...string) error {
	rootCmd := newRoot()
	rootCmd.PersistentFlags().StringVarP(&sourcesPath, "path", "p", "", "REQUIRED. Path to the sources folder")
	rootCmd.PersistentFlags().BoolVarP(&diag, "diag", "d", false, "Show application diagnostic after run")

	rootCmd.AddCommand(newInfo())
	rootCmd.AddCommand(newLostFiles())
	rootCmd.AddCommand(newLostProjects())
	rootCmd.AddCommand(newNuget())
	rootCmd.AddCommand(newVersion())
	rootCmd.AddCommand(newValidate())

	if args != nil && len(args) > 0 {
		rootCmd.SetArgs(args)
	}

	start := time.Now()
	err := rootCmd.Execute()
	elapsed := time.Since(start)

	if diag {
		printMemUsage(appPrinter)
		appPrinter.cprint("<gray>Working time:</> <green>%v</>\n", elapsed)
	}

	return err
}

func init() {
	appPrinter = newPrinter(os.Stdout)
	appFileSystem = afero.NewOsFs()
	cobra.MousetrapHelpText = ""
}
