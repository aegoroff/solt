package cmd

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"solt/internal/sys"
	"solt/msvc"
)

type lostFilesOpts struct {
	removeLost  bool
	searchAll   bool
	filter      string
	sourcesPath string
	p           printer
}

func newLostFiles(c conf) *cobra.Command {
	opts := lostFilesOpts{p: c.prn()}
	var cmd = &cobra.Command{
		Use:     "lf",
		Aliases: []string{"lostfiles"},
		Short:   "Find lost files in the folder specified",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.sourcesPath = *c.globals().sourcesPath
			return executeLostFilesCommand(opts, c.fs())
		},
	}

	cmd.Flags().StringVarP(&opts.filter, "file", "f", ".cs", "Lost files filter extension. If not set .cs extension used")
	cmd.Flags().BoolVarP(&opts.removeLost, "remove", "r", false, "Remove lost files")
	cmd.Flags().BoolVarP(&opts.searchAll, "all", "a", false, "Search all lost files including that have links to but not exists in file system")

	return cmd
}

func executeLostFilesCommand(opts lostFilesOpts, fs afero.Fs) error {
	filecollect := newFileCollector(opts.filter)
	foldcollect := newFoldersCollector()

	foldersTree := msvc.ReadSolutionDir(opts.sourcesPath, fs, filecollect, foldcollect)

	projects := msvc.SelectProjects(foldersTree)

	filer := sys.NewFiler(fs, opts.p.writer())
	logic := newLostFilesLogic(opts.searchAll, filecollect.files, foldcollect.folders, filer)
	logic.initialize(projects)

	lostFiles, err := logic.find()

	if err != nil {
		return err
	}

	s := newScreener(opts.p)
	s.writeSlice(lostFiles)

	if len(logic.unexistFiles) > 0 {
		opts.p.cprint("\n<red>These files included into projects but not exist in the file system.</>\n")

		s.writeMap(logic.unexistFiles, "Project")
	}

	if opts.removeLost {
		logic.remove(lostFiles)
	}

	return nil
}
