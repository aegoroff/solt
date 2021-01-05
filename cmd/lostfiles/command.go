package lostfiles

import (
	"github.com/spf13/cobra"
	"solt/cmd/api"
	"solt/internal/sys"
	"solt/msvc"
)

type lostFilesCommand struct {
	*api.BaseCommand
	removeLost bool
	searchAll  bool
	filter     string
}

// New creates new command that does lost files search
func New(c *api.Conf) *cobra.Command {
	var removeLost bool
	var searchAll bool
	var filter string

	cc := api.NewCobraCreator(c, func() api.Executor {
		exe := &lostFilesCommand{
			BaseCommand: api.NewBaseCmd(c),
			removeLost:  removeLost,
			searchAll:   searchAll,
			filter:      filter,
		}
		return api.NewExecutorShowHelp(exe, c)
	})

	cmd := cc.NewCommand("lf", "lostfiles", "Find lost files in the folder specified")

	cmd.Flags().StringVarP(&filter, "file", "f", ".cs", "Lost files filter extension.")
	cmd.Flags().BoolVarP(&removeLost, "remove", "r", false, "Remove lost files")
	cmd.Flags().BoolVarP(&searchAll, "all", "a", false, "Search all lost files including that have links to but not exists in file system")

	return cmd
}

func (c *lostFilesCommand) Execute(*cobra.Command) error {
	foundFiles := newFileCollector(c.filter)
	ignoredFolders := newIgnoredFoldersCollector()

	foldersTree := msvc.ReadSolutionDir(c.SourcesPath(), c.Fs(), foundFiles, ignoredFolders)

	projects := msvc.SelectProjects(foldersTree)

	exister := api.NewExister(c.Fs(), c.Writer())
	logic := newLostFilesLogic(c.searchAll, ignoredFolders.folders, exister)
	logic.initialize(projects)

	lf, err := newFinder(logic.includedFiles, logic.excludeFolders.Items())
	if err != nil {
		// return nil so as not to confuse user if no project found and it's normal case
		return nil
	}

	lostFiles := lf.find(foundFiles.files)

	s := api.NewScreener(c.Prn())
	s.WriteSlice(lostFiles)

	title := "<red>These files included into projects but not exist in the file system.</>"
	exister.Print(c.Prn(), title, "Project")

	if c.removeLost {
		filer := sys.NewFiler(c.Fs(), c.Writer())
		filer.Remove(lostFiles)
	}

	return nil
}
