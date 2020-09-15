package cmd

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"solt/msvc"
)

type lostFilesOpts struct {
	removeLost bool
	searchAll  bool
	filter     string
}

func newLostFiles() *cobra.Command {
	opts := lostFilesOpts{}
	var cmd = &cobra.Command{
		Use:     "lf",
		Aliases: []string{"lostfiles"},
		Short:   "Find lost files in the folder specified",
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeLostFilesCommand(opts, appFileSystem)
		},
	}

	cmd.Flags().StringVarP(&opts.filter, "file", "f", ".cs", "Lost files filter extension. If not set .cs extension used")
	cmd.Flags().BoolVarP(&opts.removeLost, "remove", "r", false, "Remove lost files")
	cmd.Flags().BoolVarP(&opts.searchAll, "all", "a", false, "Search all lost files including that have links to but not exists in file system")

	return cmd
}

func executeLostFilesCommand(opts lostFilesOpts, fs afero.Fs) error {
	lh := newLostFilesHandler(opts.filter, opts.searchAll, fs)

	foldersTree := msvc.ReadSolutionDir(sourcesPath, fs, lh)

	projects := msvc.SelectProjects(foldersTree)

	lh.projectHandler(projects)

	lostFiles, err := lh.findLostFiles()

	if err != nil {
		return err
	}

	sortAndOutput(appPrinter, lostFiles)

	if len(lh.unexistFiles) > 0 {
		appPrinter.cprint("\n<red>These files included into projects but not exist in the file system.</>\n")

		outputSortedMap(appPrinter, lh.unexistFiles, "Project")
	}

	if opts.removeLost {
		lh.removeLostFiles(lostFiles)
	}

	return nil
}
