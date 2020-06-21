package cmd

import (
	"github.com/gookit/color"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"solt/msvc"
)

type lostFilesCmd struct {
	removeLost bool
	searchAll  bool
	filter     string
}

func newLostFiles() *cobra.Command {
	opts := lostFilesCmd{}
	var cmd = &cobra.Command{
		Use:     "lf",
		Aliases: []string{"lostfiles"},
		Short:   "Find lost files in the folder specified",
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeLostFilesCommand(opts.filter, opts.removeLost, opts.searchAll, appFileSystem)
		},
	}

	cmd.Flags().StringVarP(&opts.filter, "file", "f", ".cs", "Lost files filter extension. If not set .cs extension used")
	cmd.Flags().BoolVarP(&opts.removeLost, "remove", "r", false, "Remove lost files")
	cmd.Flags().BoolVarP(&opts.searchAll, "all", "a", false, "Search all lost files including that have links to but not exists in file system")

	return cmd
}

func executeLostFilesCommand(lostFilesFilter string, removeLostFiles bool, nonExist bool, fs afero.Fs) error {
	lh := newLostFilesHandler(lostFilesFilter, nonExist, fs)

	foldersTree := msvc.ReadSolutionDir(sourcesPath, fs, lh)

	projects := msvc.SelectProjects(foldersTree)

	lh.projectHandler(projects)

	lostFiles, err := lh.findLostFiles()

	if err != nil {
		return err
	}

	sortAndOutput(appWriter, lostFiles)

	if len(lh.unexistFiles) > 0 {
		color.Fprintf(appWriter, "\n<red>These files included into projects but not exist in the file system.</>\n")

		outputSortedMap(appWriter, lh.unexistFiles, "Project")
	}

	if removeLostFiles {
		lh.removeLostFiles(lostFiles)
	}

	return nil
}
