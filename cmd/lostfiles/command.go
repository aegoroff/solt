package lostfiles

import (
	"github.com/spf13/cobra"
	"solt/cmd/fw"
	"solt/internal/sys"
	"solt/msvc"
)

type lostFilesCommand struct {
	*fw.BaseCommand
	removeLost bool
	searchAll  bool
	filter     string
}

// New creates new command that does lost files search
func New(c *fw.Conf) *cobra.Command {
	var removeLost bool
	var searchAll bool
	var filter string

	cc := fw.NewCobraCreator(c, func() fw.Executor {
		exe := &lostFilesCommand{
			BaseCommand: fw.NewBaseCmd(c),
			removeLost:  removeLost,
			searchAll:   searchAll,
			filter:      filter,
		}
		return fw.NewExecutorShowHelp(exe, c)
	})

	cmd := cc.NewCommand("lf", "lostfiles", "Find lost files in the folder specified")

	cmd.Flags().StringVarP(&filter, "file", "f", ".cs", "Lost files filter extension.")
	cmd.Flags().BoolVarP(&removeLost, "remove", "r", false, "Remove lost files")
	cmd.Flags().BoolVarP(&searchAll, "all", "a", false, "Search all lost files including that have links to but not exists in file system")

	return cmd
}

func (c *lostFilesCommand) Execute(*cobra.Command) error {
	exist := newExister(c.searchAll, c.Fs(), c.Writer())
	incl := newIncluder(exist)

	collect := newCollector(c.filter)
	skip := newSkipper()

	foldersTree := msvc.ReadSolutionDir(c.SourcesPath(), c.Fs(), collect, skip)
	projects := msvc.SelectProjects(foldersTree)
	enumerate(projects, skip.fromProject, incl.fromProject)

	lf, err := newFinder(incl.files(), skip.folders())
	if err != nil {
		// return nil so as not to confuse user if no project found and it's normal case
		return nil
	}

	lost := lf.find(collect.files)

	c.print(lost)

	exist.print(c.Prn())

	c.remove(lost)

	return nil
}

func (c *lostFilesCommand) print(lost []string) {
	s := fw.NewScreener(c.Prn())
	s.WriteSlice(lost)
}

func (c *lostFilesCommand) remove(lost []string) {
	if c.removeLost {
		filer := sys.NewFiler(c.Fs(), c.Writer())
		filer.Remove(lost)
	}
}
