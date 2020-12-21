package cmd

import (
	"github.com/dustin/go-humanize"
	"github.com/spf13/afero"
	"io"
	"runtime"
	"time"

	"github.com/spf13/cobra"
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
	var diag bool

	rootCmd.PersistentFlags().StringVarP(&sourcesPath, "path", "p", "", "REQUIRED. Path to the sources folder")
	rootCmd.PersistentFlags().BoolVarP(&diag, "diag", "d", false, "Show application diagnostic after run")

	conf := newAppConf(fs, p, &sourcesPath)

	rootCmd.AddCommand(newInfo(conf))
	rootCmd.AddCommand(newLostFiles(conf))
	rootCmd.AddCommand(newLostProjects(conf))
	rootCmd.AddCommand(newNuget(conf))
	rootCmd.AddCommand(newVersion(conf))
	rootCmd.AddCommand(newValidate(conf))

	if args != nil && len(args) > 0 {
		rootCmd.SetArgs(args)
	}

	start := time.Now()
	err := rootCmd.Execute()
	elapsed := time.Since(start)

	if diag {
		printMemUsage(conf.prn())
		conf.prn().cprint("<gray>Working time:</> <green>%v</>\n", elapsed)
	}

	return err
}

func init() {
	cobra.MousetrapHelpText = ""
}

// printMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func printMemUsage(p printer) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	p.cprint("\n<gray>Alloc =</> <green>%s</>", humanize.IBytes(m.Alloc))
	p.cprint("\t<gray>TotalAlloc =</> <green>%s</>", humanize.IBytes(m.TotalAlloc))
	p.cprint("\t<gray>Sys =</> <green>%s</>", humanize.IBytes(m.Sys))
	p.cprint("\t<gray>NumGC =</> <green>%v</>\n", m.NumGC)
}
