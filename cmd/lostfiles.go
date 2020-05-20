package cmd

import (
	"fmt"
	"github.com/aegoroff/godatastruct/collections"
	"github.com/aegoroff/godatastruct/rbtree"
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

		lostFiles, unexistFiles := findLostFiles(appFileSystem, foldersTree, packagesFolders.Items(), foundFiles)

		sortAndOutput(appWriter, lostFiles)

		if len(unexistFiles) > 0 {
			fmt.Fprintf(appWriter, "\nThese files included into projects but not exist in the file system.\n")
		}

		outputSortedMap(appWriter, unexistFiles, "Project")

		if removeLostFiles {
			for _, f := range lostFiles {
				err = appFileSystem.Remove(f)
				if err != nil {
					log.Printf("%v\n", err)
				} else {
					fmt.Fprintf(appWriter, "File: %s removed sucessfully.\n", f)
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

func findLostFiles(fs afero.Fs, foldersTree *rbtree.RbTree, additionalFoldersToExclude []string, foundFiles []string) ([]string, map[string][]string) {
	includedFiles, excludedFolders, unexistFiles := createIncludedFilesAndExcludedFolders(foldersTree, fs)
	excludedFolders = append(excludedFolders, additionalFoldersToExclude...)

	exmach, err := createAhoCorasickMachine(excludedFolders)
	if err != nil {
		fmt.Println(err)
		return []string{}, unexistFiles
	}
	var result []string
	for _, file := range foundFiles {
		if !includedFiles.Contains(strings.ToUpper(file)) && !Match(exmach, file) {
			result = append(result, file)
		}
	}

	return result, unexistFiles
}

func createIncludedFilesAndExcludedFolders(foldersTree *rbtree.RbTree, fs afero.Fs) (collections.StringHashSet, []string, map[string][]string) {
	var excludeFolders []string
	unexistFiles := make(map[string][]string)
	var includedFiles = make(collections.StringHashSet)

	foldersTree.Ascend(func(key *rbtree.Comparable) bool {
		folder := (*key).(*folder)
		content := folder.content

		if len(content.projects) == 0 {
			return true
		}

		for _, prj := range content.projects {
			// Add project base + exclude subfolder into exclude folders list
			for _, s := range subfolderToExclude {
				sub := filepath.Join(folder.path, s)
				excludeFolders = append(excludeFolders, sub)
			}

			// Exclude output paths too
			if prj.project.OutputPaths != nil {
				for _, out := range prj.project.OutputPaths {
					sub := filepath.Join(folder.path, out)
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
		}

		return true
	})

	return includedFiles, excludeFolders, unexistFiles
}
