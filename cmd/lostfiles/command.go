package lostfiles

import (
	"github.com/spf13/cobra"
	"solt/cmd/api"
	"solt/internal/sys"
	"solt/msvc"
)

type lostFilesCommand struct {
	api.BaseCommand
	removeLost bool
	searchAll  bool
	filter     string
}

func New(c *api.Conf) *cobra.Command {
	var removeLost bool
	var searchAll bool
	var filter string

	cc := api.NewCobraCreator(c, func() api.Executor {
		return &lostFilesCommand{
			BaseCommand: api.NewBaseCmd(c),
			removeLost:  removeLost,
			searchAll:   searchAll,
			filter:      filter,
		}
	})

	cmd := cc.NewCommand("lf", "lostfiles", "Find lost files in the folder specified")

	cmd.Flags().StringVarP(&filter, "file", "f", ".cs", "Lost files filter extension.")
	cmd.Flags().BoolVarP(&removeLost, "remove", "r", false, "Remove lost files")
	cmd.Flags().BoolVarP(&searchAll, "all", "a", false, "Search all lost files including that have links to but not exists in file system")

	return cmd
}

func (c *lostFilesCommand) Execute() error {
	filecollect := newFileCollector(c.filter)
	foldcollect := newFoldersCollector()

	foldersTree := msvc.ReadSolutionDir(c.SourcesPath(), c.Fs(), filecollect, foldcollect)

	projects := msvc.SelectProjects(foldersTree)

	filer := sys.NewFiler(c.Fs(), c.Prn().Writer())
	logic := newLostFilesLogic(c.searchAll, filecollect.files, foldcollect.folders, filer)
	err := logic.initialize(projects)

	if err != nil {
		return err
	}

	lostFiles := logic.find()

	s := api.NewScreener(c.Prn())
	s.WriteSlice(lostFiles)

	if len(logic.unexistFiles) > 0 {
		c.Prn().Cprint("\n<red>These files included into projects but not exist in the file system.</>\n")

		s.WriteMap(logic.unexistFiles, "Project")
	}

	if c.removeLost {
		logic.remove(lostFiles)
	}

	return nil
}
