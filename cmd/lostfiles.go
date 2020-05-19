package cmd

import (
	"fmt"
	"github.com/aegoroff/godatastruct/collections"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

var subfolderToExclude = []string{
	"obj",
}

const filterParamName = "file"

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

		lostFiles, unexistFiles := findLostFiles(foldersTree, appFileSystem, packagesFolders, foundFiles)

		sortAndOutput(appWriter, lostFiles)

		if len(unexistFiles) > 0 {
			fmt.Fprintf(appWriter, "\nThese files included into projects but not exist in the file system.\n")
		}

		outputSortedMap(appWriter, unexistFiles, "Project")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(lostfilesCmd)
	lostfilesCmd.Flags().StringP(filterParamName, "f", ".cs", "Lost files filter extension. If not set .cs extension used")
}

func findLostFiles(foldersTree *rbtree.RbTree, fs afero.Fs, packagesFolders collections.StringHashSet, foundFiles []string) ([]string, map[string][]string) {
	includedFiles, excludedFolders, unexistFiles := createIncludedFilesAndExcludedFolders(foldersTree, fs)
	excludedFolders = append(excludedFolders, packagesFolders.Items()...)

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
			project := prj.file

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

			// Add compiles, contents and nones into included files map
			filesIncluded := getFilesIncludedIntoProject(prj)
			for _, f := range filesIncluded {
				includedFiles.Add(strings.ToUpper(f))
				if _, err := fs.Stat(f); os.IsNotExist(err) {
					if found, ok := unexistFiles[project]; ok {
						found = append(found, f)
						unexistFiles[project] = found
					} else {
						unexistFiles[project] = []string{f}
					}
				}
			}
		}

		return true
	})

	return includedFiles, excludeFolders, unexistFiles
}
