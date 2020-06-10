package cmd

import (
	"fmt"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"solt/msvc"
)

const filterParamName = "file"
const removeParamName = "remove"
const allParamName = "all"

func newLostFiles() *cobra.Command {
	var lostfilesCmd = &cobra.Command{
		Use:     "lf",
		Aliases: []string{"lostfiles"},
		Short:   "Find lost files in the folder specified",
		RunE: func(cmd *cobra.Command, args []string) error {
			lostFilesFilter, err := cmd.Flags().GetString(filterParamName)

			if err != nil {
				return err
			}

			removeLostFiles, err := cmd.Flags().GetBool(removeParamName)

			if err != nil {
				return err
			}

			all, err := cmd.Flags().GetBool(allParamName)

			if err != nil {
				return err
			}

			return executeLostFilesCommand(lostFilesFilter, removeLostFiles, all, appFileSystem)
		},
	}

	lostfilesCmd.Flags().StringP(filterParamName, "f", ".cs", "Lost files filter extension. If not set .cs extension used")
	lostfilesCmd.Flags().BoolP(removeParamName, "r", false, "Remove lost files")
	lostfilesCmd.Flags().BoolP(allParamName, "a", false, "Search all lost files including that have links to but not exists in file system")

	return lostfilesCmd
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
		_, _ = fmt.Fprintf(appWriter, "\nThese files included into projects but not exist in the file system.\n")

		outputSortedMap(appWriter, lh.unexistFiles, "Project")
	}

	if removeLostFiles {
		lh.removeLostFiles(lostFiles)
	}

	return nil
}
