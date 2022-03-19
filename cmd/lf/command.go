package lf

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"io"
	"solt/internal/fw"
	"solt/internal/sys"
	"solt/internal/ux"
	"solt/msvc"
)

type lostFilesCommand struct {
	*fw.BaseCommand
	remover fw.Remover
	exister fw.Exister
	filter  string
}

// New creates new command that does lost files search
func New(c *fw.Conf) *cobra.Command {
	var removeLost bool
	var searchAll bool
	var filter string

	cc := fw.NewCobraCreator(c, func() fw.Executor {
		exe := &lostFilesCommand{
			BaseCommand: fw.NewBaseCmd(c),
			remover:     newRemover(c.Fs(), c.W(), removeLost),
			exister:     newExister(c.Fs(), c.W(), searchAll),
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

func newExister(fs afero.Fs, w io.Writer, real bool) fw.Exister {
	if real {
		return fw.NewExister(fs, w)
	}
	return &nullExister{}
}

func newRemover(fs afero.Fs, w io.Writer, real bool) fw.Remover {
	if real {
		return sys.NewFiler(fs, w)
	}
	return &nullRemover{}
}

func (c *lostFilesCommand) Execute(*cobra.Command) error {
	collected := newCollector(c.filter)
	skip := newSkipper()

	foldersTree := msvc.ReadSolutionDir(c.SourcesPath(), c.Fs(), collected, skip)
	projects := msvc.SelectProjects(foldersTree)

	incl := fw.NewIncluder(c.exister)

	for _, p := range projects {
		skip.fromProject(p)
		incl.From(p)
	}

	filesIn := incl.Includes()
	lf := newFinder(filesIn, skip.folders())
	lost := lf.find(collected.files)

	c.print(lost)

	c.remover.Remove(lost)

	tt := &totals{
		projects: int64(len(projects)),
		missing:  c.exister.MissingCount(),
		included: int64(len(filesIn)),
		lost:     int64(len(lost)),
		found:    int64(len(collected.files)),
	}
	c.Prn().Println()
	c.Total(tt)
	return nil
}

func (c *lostFilesCommand) print(lost []string) {
	s := ux.NewScreener(c.Prn())
	s.WriteSlice(lost)

	title := "<red>These files included into projects but not exist in the file system.</>"
	c.exister.Print(c.Prn(), title, "Project")
}
