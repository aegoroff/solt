package cmd

import (
	"fmt"
	"github.com/aegoroff/godatastruct/collections"
	"github.com/aegoroff/godatastruct/rbtree"
	goahocorasick "github.com/anknown/ahocorasick"
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

// lostfilesCmd represents the lostfiles command
var lostfilesCmd = &cobra.Command{
	Use:     "lostfiles",
	Aliases: []string{"lf"},
	Short:   "Find lost files in the folder specified",
	RunE: func(cmd *cobra.Command, args []string) error {
		var foundFiles []string
		var packagesFolders = make(collections.StringHashSet)

		lostFilesFilter, err := cmd.Flags().GetString(filterParamName)

		if err != nil {
			return err
		}

		removeLostFiles, err := cmd.Flags().GetBool(removeParamName)

		if err != nil {
			return err
		}

		foldersTree := readProjectDir(sourcesPath, appFileSystem, func(we *walkEntry) {
			// Add file to filtered files slice
			ext := strings.ToLower(filepath.Ext(we.Name))
			if ext == lostFilesFilter {
				fp := filepath.Join(we.Parent, we.Name)
				foundFiles = append(foundFiles, fp)
			}

			if ext == solutionFileExt {
				ppath := filepath.Join(we.Parent, "packages")
				packagesFolders.Add(ppath)
			}
		})

		includedFiles, excludedFolders, unexistFiles := createIncludedFilesAndExcludedFolders(foldersTree, appFileSystem)
		excludedFolders = append(excludedFolders, packagesFolders.Items()...)

		exmach, err := createAhoCorasickMachine(excludedFolders)
		if err != nil {
			return err
		}

		lostFiles := findLostFiles(exmach, foundFiles, includedFiles)

		sortAndOutput(appWriter, lostFiles)

		if len(unexistFiles) > 0 {
			_, _ = fmt.Fprintf(appWriter, "\nThese files included into projects but not exist in the file system.\n")
		}

		outputSortedMap(appWriter, unexistFiles, "Project")

		if removeLostFiles {
			for _, f := range lostFiles {
				err = appFileSystem.Remove(f)
				if err != nil {
					log.Printf("%v\n", err)
				} else {
					_, _ = fmt.Fprintf(appWriter, "File: %s removed sucessfully.\n", f)
				}
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(lostfilesCmd)
	lostfilesCmd.Flags().StringP(filterParamName, "f", ".cs", "Lost files filter extension. If not set .cs extension used")
	lostfilesCmd.Flags().BoolP(removeParamName, "r", false, "Remove lost files")
}

func findLostFiles(exm *goahocorasick.Machine, foundFiles []string, includedFiles collections.StringHashSet) []string {
	var result []string
	for _, file := range foundFiles {
		if !includedFiles.Contains(strings.ToUpper(file)) && !Match(exm, file) {
			result = append(result, file)
		}
	}

	return result
}

func createIncludedFilesAndExcludedFolders(foldersTree *rbtree.RbTree, fs afero.Fs) (collections.StringHashSet, []string, map[string][]string) {
	var excludeFolders []string
	unexistFiles := make(map[string][]string)
	var includedFiles = make(collections.StringHashSet)

	walkProjects(foldersTree, func(prj *msbuildProject, fold *folder) {
		// Add project base + exclude subfolder into exclude folders list
		for _, s := range subfolderToExclude {
			sub := filepath.Join(fold.path, s)
			excludeFolders = append(excludeFolders, sub)
		}

		// Exclude output paths too
		if prj.project.OutputPaths != nil {
			for _, out := range prj.project.OutputPaths {
				sub := filepath.Join(fold.path, out)
				excludeFolders = append(excludeFolders, sub)
			}
		}

		// In case of SDK projects all files inside project folder are considered included
		if prj.project.isSdkProject() {
			excludeFolders = append(excludeFolders, filepath.Dir(prj.path))
		}

		// Add compiles, contents and nones into included files map
		filesIncluded := getFilesIncludedIntoProject(prj)
		for _, f := range filesIncluded {
			includedFiles.Add(strings.ToUpper(f))
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

	return includedFiles, excludeFolders, unexistFiles
}
