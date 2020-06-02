package cmd

import (
	"fmt"
	"github.com/aegoroff/godatastruct/collections"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var subfolderToExclude = []string{
	"obj",
}

const filterParamName = "file"
const removeParamName = "remove"
const onlyLostParamName = "onlylost"

// lostfilesCmd represents the lostfiles command
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

		onlyLost, err := cmd.Flags().GetBool(onlyLostParamName)

		if err != nil {
			return err
		}

		return executeLostFilesCommand(lostFilesFilter, removeLostFiles, onlyLost, appFileSystem)
	},
}

func init() {
	rootCmd.AddCommand(lostfilesCmd)
	lostfilesCmd.Flags().StringP(filterParamName, "f", ".cs", "Lost files filter extension. If not set .cs extension used")
	lostfilesCmd.Flags().BoolP(removeParamName, "r", false, "Remove lost files")
	lostfilesCmd.Flags().BoolP(onlyLostParamName, "l", false, "Show only lost files. Don't show unexist files. If not set all shown")
}

func executeLostFilesCommand(lostFilesFilter string, removeLostFiles bool, onlyLost bool, fs afero.Fs) error {
	var foundFiles []string
	var excludeFolders = make(collections.StringHashSet)
	ef := normalize(lostFilesFilter)
	sln := normalize(solutionFileExt)
	foldersTree := readProjectDir(sourcesPath, fs, func(path string) {
		// Add file to filtered files slice
		ext := normalize(filepath.Ext(path))
		if ext == ef {
			foundFiles = append(foundFiles, path)
		}

		if ext == sln {
			dir, _ := filepath.Split(path)
			ppath := filepath.Join(dir, "packages")
			excludeFolders.Add(ppath)
		}
	})

	unexistFiles := make(map[string][]string)
	var includedFiles = make(collections.StringHashSet)

	walkProjects(foldersTree, func(prj *msbuildProject, fold *folder) {
		// Add project base + exclude subfolder into exclude folders list
		for _, s := range subfolderToExclude {
			sub := filepath.Join(fold.path, s)
			excludeFolders.Add(sub)
		}

		// Exclude output paths too
		if prj.project.OutputPaths != nil {
			for _, out := range prj.project.OutputPaths {
				sub := filepath.Join(fold.path, out)
				excludeFolders.Add(sub)
			}
		}

		// In case of SDK projects all files inside project folder are considered included
		if prj.project.isSdkProject() {
			excludeFolders.Add(filepath.Dir(prj.path))
		}

		// Add compiles, contents and nones into included files map
		filesIncluded := getFilesIncludedIntoProject(prj)
		for _, f := range filesIncluded {
			normalized := normalize(f)
			includedFiles.Add(normalized)
			if _, err := fs.Stat(f); os.IsNotExist(err) {
				if found, ok := unexistFiles[prj.path]; ok {
					found = append(found, f)
					unexistFiles[prj.path] = found
				} else {
					unexistFiles[prj.path] = []string{f}
				}
			}
		}
	})

	lostFiles, err := findLostFiles(excludeFolders, foundFiles, includedFiles)

	if err != nil {
		return err
	}

	sortAndOutput(appWriter, lostFiles)

	if !onlyLost {
		if len(unexistFiles) > 0 {
			_, _ = fmt.Fprintf(appWriter, "\nThese files included into projects but not exist in the file system.\n")
		}

		outputSortedMap(appWriter, unexistFiles, "Project")
	}

	if removeLostFiles {
		removeLostfiles(lostFiles, fs)
	}

	if showMemUsage {
		printMemUsage(appWriter)
	}
	return nil
}

func findLostFiles(excludeFolders collections.StringHashSet, foundFiles []string, includedFiles collections.StringHashSet) ([]string, error) {
	exm, err := createAhoCorasickMachine(excludeFolders.ItemsDecorated(normalize))
	if err != nil {
		return nil, err
	}

	var result []string
	for _, file := range foundFiles {
		normalized := normalize(file)
		if !includedFiles.Contains(normalized) && !Match(exm, normalized) {
			result = append(result, file)
		}
	}

	return result, err
}

func normalize(s string) string {
	return strings.ToUpper(s)
}

func removeLostfiles(lostFiles []string, fs afero.Fs) {
	for _, f := range lostFiles {
		err := fs.Remove(f)
		if err != nil {
			log.Printf("%v\n", err)
		} else {
			_, _ = fmt.Fprintf(appWriter, "File: %s removed successfully.\n", f)
		}
	}
}
