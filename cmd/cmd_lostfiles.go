package cmd

import (
	"github.com/spf13/cobra"
	"solt/internal/sys"
	"solt/msvc"
)

type lostFilesCommand struct {
	baseCommand
	removeLost bool
	searchAll  bool
	filter     string
}

func newLostFiles(c *conf) *cobra.Command {
	var removeLost bool
	var searchAll bool
	var filter string

	cc := cobraCreator{
		createCmd: func() executor {
			return &lostFilesCommand{
				baseCommand: newBaseCmd(c),
				removeLost:  removeLost,
				searchAll:   searchAll,
				filter:      filter,
			}
		},
		c: c,
	}

	cmd := cc.newCobraCommand("lf", "lostfiles", "Find lost files in the folder specified")

	cmd.Flags().StringVarP(&filter, "file", "f", ".cs", "Lost files filter extension.")
	cmd.Flags().BoolVarP(&removeLost, "remove", "r", false, "Remove lost files")
	cmd.Flags().BoolVarP(&searchAll, "all", "a", false, "Search all lost files including that have links to but not exists in file system")

	return cmd
}

func (c *lostFilesCommand) execute() error {
	filecollect := newFileCollector(c.filter)
	foldcollect := newFoldersCollector()

	foldersTree := msvc.ReadSolutionDir(c.sourcesPath, c.fs, filecollect, foldcollect)

	projects := msvc.SelectProjects(foldersTree)

	filer := sys.NewFiler(c.fs, c.prn.writer())
	logic := newLostFilesLogic(c.searchAll, filecollect.files, foldcollect.folders, filer)
	err := logic.initialize(projects)

	if err != nil {
		return err
	}

	lostFiles := logic.find()

	s := newScreener(c.prn)
	s.writeSlice(lostFiles)

	if len(logic.unexistFiles) > 0 {
		c.prn.cprint("\n<red>These files included into projects but not exist in the file system.</>\n")

		s.writeMap(logic.unexistFiles, "Project")
	}

	if c.removeLost {
		logic.remove(lostFiles)
	}

	return nil
}
