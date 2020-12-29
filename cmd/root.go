package cmd

import (
	"github.com/aegoroff/dirstat/scan"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"io"
	"log"
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
	rootCmd := newRoot()

	var sourcesPath string
	var cpuprofile string
	var memprofile string
	var resultfile string
	var diag bool

	rootCmd.PersistentFlags().StringVarP(&sourcesPath, "path", "p", "", "REQUIRED. Path to the sources folder")

	const profileTrail = "If not set profiling not started. Correct file path should be set here"

	const cDescription = "Runs CPU profiling if --diag option set. " + profileTrail
	rootCmd.PersistentFlags().StringVarP(&cpuprofile, "cpuprofile", "", "", cDescription)

	const mDescription = "Runs memory profiling if --diag option set. " + profileTrail
	rootCmd.PersistentFlags().StringVarP(&memprofile, "memprofile", "", "", mDescription)

	const fDescription = "Write results into file. Specify path to output file using this option"
	rootCmd.PersistentFlags().StringVarP(&resultfile, "output", "o", "", fDescription)

	rootCmd.PersistentFlags().BoolVarP(&diag, "diag", "d", false, "Show application diagnostic after run")

	env := newWriteFileEnvironment(&resultfile, fs, pe)
	defer env.close()

	c := api.NewConf(fs, env, &sourcesPath, &cpuprofile, &memprofile, &diag)

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

type fileEnvironment struct {
	path *string
	fs   afero.Fs
	pe   api.PrintEnvironment
	file afero.File
}

func newWriteFileEnvironment(path *string, fs afero.Fs, defaultpe api.PrintEnvironment) *fileEnvironment {
	pe := &fileEnvironment{
		path: path,
		fs:   fs,
		pe:   defaultpe,
	}
	return pe
}

func (e *fileEnvironment) create(path *string, fs afero.Fs) error {
	f, err := fs.Create(*path)
	if err != nil {
		return err
	}

	e.pe = api.NewStringEnvironment(f)
	e.file = f

	return nil
}

func (e *fileEnvironment) close() {
	scan.Close(e.file)
}

func (e *fileEnvironment) NewPrinter() api.Printer {
	if *e.path == "" {
		return e.pe.NewPrinter()
	}
	err := e.create(e.path, e.fs)
	if err != nil {
		log.Println(err)
		return e.pe.NewPrinter()
	}
	return api.NewPrinter(e)
}

func (e *fileEnvironment) PrintFunc(w io.Writer, format string, a ...interface{}) {
	e.pe.PrintFunc(w, format, a...)
}

func (e *fileEnvironment) Writer() io.Writer {
	return e.pe.Writer()
}

func init() {
	cobra.MousetrapHelpText = ""
}
