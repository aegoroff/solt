package lostfiles

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"io"
	"solt/cmd/fw"
	"solt/internal/sys"
	"solt/internal/ux"
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
	collect := newCollector(c.filter)
	skip := newSkipper()

	foldersTree := msvc.ReadSolutionDir(c.SourcesPath(), c.Fs(), collect, skip)
	projects := msvc.SelectProjects(foldersTree)

	exist := c.newExister(c.Fs(), c.Writer())
	incl := fw.NewIncluder(exist)

	for _, p := range projects {
		skip.fromProject(p)
		incl.From(p)
	}

	lf := newFinder(incl.Includes(), skip.folders())
	lost := lf.find(collect.files)

	c.print(lost)

	title := "<red>These files included into projects but not exist in the file system.</>"
	exist.Print(c.Prn(), title, "Project")

	c.removeIfRequested(lost)

	return nil
}

func (c *lostFilesCommand) newExister(fs afero.Fs, w io.Writer) fw.Exister {
	if c.searchAll {
		return fw.NewExister(fs, w)
	}
	return fw.NewNullExister()
}

func (c *lostFilesCommand) print(lost []string) {
	s := ux.NewScreener(c.Prn())
	s.WriteSlice(lost)
}

func (c *lostFilesCommand) removeIfRequested(lost []string) {
	if c.removeLost {
		filer := sys.NewFiler(c.Fs(), c.Writer())
		filer.Remove(lost)
	}
}
