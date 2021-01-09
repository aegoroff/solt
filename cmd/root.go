package cmd

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"solt/cmd/fw"
	"solt/cmd/info"
	"solt/cmd/lostfiles"
	"solt/cmd/lostprojects"
	"solt/cmd/nuget"
	"solt/cmd/validate"
	"solt/internal/out"
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
func Execute(fs afero.Fs, pe out.PrintEnvironment, args ...string) error {
	rootCmd := newRoot()

	var resultfile string

	d := fw.Diag{}

	const profileTrail = "If not set profiling not started. Correct file path should be set here"

	const cDescription = "Runs CPU profiling if --diag option set. " + profileTrail
	rootCmd.PersistentFlags().StringVarP(&d.CPU, "cpuprofile", "", "", cDescription)

	const mDescription = "Runs memory profiling if --diag option set. " + profileTrail
	rootCmd.PersistentFlags().StringVarP(&d.Memory, "memprofile", "", "", mDescription)

	const fDescription = "Write results into file. Specify path to output file using this option"
	rootCmd.PersistentFlags().StringVarP(&resultfile, "output", "o", "", fDescription)

	rootCmd.PersistentFlags().BoolVarP(&d.Enable, "diag", "d", false, "Show application diagnostic after run")

	env := out.NewWriteFileEnvironment(&resultfile, fs, pe)

	c := fw.NewConf(fs, env, &d)

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
