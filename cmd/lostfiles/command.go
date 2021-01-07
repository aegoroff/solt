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
	foundFiles := newFileCollector(c.filter)
	skip := newSkipper()

	foldersTree := msvc.ReadSolutionDir(c.SourcesPath(), c.Fs(), foundFiles, skip)

	projects := msvc.SelectProjects(foldersTree)

	exist := fw.NewExister(c.Fs(), c.Writer())
	wrapper := newExister(c.searchAll, exist)
	enum := newEnumerator(skip, wrapper)
	enum.enumerate(projects)

	lf, err := newFinder(enum.includedFiles(), skip.skipped())
	if err != nil {
		// return nil so as not to confuse user if no project found and it's normal case
		return nil
	}

	lostFiles := lf.find(foundFiles.files)

	s := fw.NewScreener(c.Prn())
	s.WriteSlice(lostFiles)

	title := "<red>These files included into projects but not exist in the file system.</>"
	exist.Print(c.Prn(), title, "Project")

	if c.removeLost {
		filer := sys.NewFiler(c.Fs(), c.Writer())
		filer.Remove(lostFiles)
	}

	return nil
}
