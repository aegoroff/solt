package cmd

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"solt/cmd/api"
	"solt/cmd/info"
	"solt/cmd/lostfiles"
	"solt/cmd/lostprojects"
	"solt/cmd/nuget"
	"solt/cmd/validate"
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
func Execute(fs afero.Fs, pe api.PrintEnvironment, args ...string) error {
	p := api.NewPrinter(pe)

	rootCmd := newRoot()

	var sourcesPath string
	var cpuprofile string
	var memprofile string
	var diag bool

	rootCmd.PersistentFlags().StringVarP(&sourcesPath, "path", "p", "", "REQUIRED. Path to the sources folder")

	const profileTrail = "If not set profiling not started. Correct file path should be set here"

	const cDescription = "Runs CPU profiling if --diag option set. " + profileTrail
	rootCmd.PersistentFlags().StringVarP(&cpuprofile, "cpuprofile", "", "", cDescription)

	const mDescription = "Runs memory profiling if --diag option set. " + profileTrail
	rootCmd.PersistentFlags().StringVarP(&memprofile, "memprofile", "", "", mDescription)

	rootCmd.PersistentFlags().BoolVarP(&diag, "diag", "d", false, "Show application diagnostic after run")

	c := api.NewConf(fs, p, &sourcesPath, &cpuprofile, &memprofile, &diag)

	rootCmd.AddCommand(info.New(c))
	rootCmd.AddCommand(lostfiles.New(c))
	rootCmd.AddCommand(lostprojects.New(c))
	rootCmd.AddCommand(nuget.New(c))
	rootCmd.AddCommand(newVersion(c))
	rootCmd.AddCommand(validate.New(c))

	if args != nil && len(args) > 0 {
		rootCmd.SetArgs(args)
	}

	return rootCmd.Execute()
}

func init() {
	cobra.MousetrapHelpText = ""
}
