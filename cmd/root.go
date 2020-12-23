package cmd

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"io"
)

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
func Execute(fs afero.Fs, w io.Writer, args ...string) error {
	p := newPrinter(w)
	return execute(fs, p, args...)
}

func execute(fs afero.Fs, p printer, args ...string) error {
	rootCmd := newRoot()

	var sourcesPath string
	var cpuprofile string
	var diag bool

	rootCmd.PersistentFlags().StringVarP(&sourcesPath, "path", "p", "", "REQUIRED. Path to the sources folder")
	const cDescr = "Runs CPU profiling if --diag option set. If not set profiling not started. Correct file path should be set here"
	rootCmd.PersistentFlags().StringVarP(&cpuprofile, "cpuprofile", "", "", cDescr)
	rootCmd.PersistentFlags().BoolVarP(&diag, "diag", "d", false, "Show application diagnostic after run")

	c := &conf{
		filesystem: fs,
		p:          p,
		sp:         &sourcesPath,
		cpu:        &cpuprofile,
		diag:       &diag,
	}

	rootCmd.AddCommand(newInfo(c))
	rootCmd.AddCommand(newLostFiles(c))
	rootCmd.AddCommand(newLostProjects(c))
	rootCmd.AddCommand(newNuget(c))
	rootCmd.AddCommand(newVersion(c))
	rootCmd.AddCommand(newValidate(c))

	if args != nil && len(args) > 0 {
		rootCmd.SetArgs(args)
	}

	return rootCmd.Execute()
}

func init() {
	cobra.MousetrapHelpText = ""
}
