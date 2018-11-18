package cmd

import (
	"fmt"
	"github.com/aegoroff/godatastruct/rbtree"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var subfolderToExclude = []string{
	"obj",
}

const filterParamName = "file"

// lostfilesCmd represents the lostfiles command
var lostfilesCmd = &cobra.Command{
	Use:   "lostfiles",
	Short: "Find lost files in the folder specified",
	Run: func(cmd *cobra.Command, args []string) {
		var foundFiles []string
		var packagesFolders = make(map[string]interface{})

		lostFilesFilter, _ := cmd.Flags().GetString(filterParamName)

		foldersTree := readProjectDir(sourcesPath, func(we *walkEntry) {
			// Add file to filtered files slice
			ext := strings.ToLower(filepath.Ext(we.Name))
			if ext == lostFilesFilter {
				fp := filepath.Join(we.Parent, we.Name)
				foundFiles = append(foundFiles, fp)
			}

			if ext == solutionFileExt {
				ppath := filepath.Join(we.Parent, "packages")
				if _, ok := packagesFolders[ppath]; !ok {
					packagesFolders[ppath] = nil
				}
			}
		})

		lostFiles, unexistFiles := findLostFiles(foldersTree, packagesFolders, foundFiles)

		sortAndOutputToStdout(lostFiles)

		if len(unexistFiles) > 0 {
			fmt.Printf("\nThese files included into projects but not exist in the file system.\n")
		}

		outputSortedMapToStdout(unexistFiles, "Project")
	},
}

func init() {
	rootCmd.AddCommand(lostfilesCmd)
	lostfilesCmd.Flags().StringP(filterParamName, "f", ".cs", "Lost files filter extension")
}

func findLostFiles(foldersTree *rbtree.RbTree, packagesFolders map[string]interface{}, foundFiles []string) ([]string, map[string][]string) {
	includedFiles, excludedFolders, unexistFiles := createIncludedFilesAndExcludedFolders(foldersTree)
	for k := range packagesFolders {
		excludedFolders = append(excludedFolders, k)
	}

	exmach := createAhoCorasickMachine(excludedFolders)
	var result []string
	for _, file := range foundFiles {
		if _, ok := includedFiles[strings.ToUpper(file)]; !ok && !Match(exmach, file) {
			result = append(result, file)
		}
	}

	return result, unexistFiles
}

func createIncludedFilesAndExcludedFolders(foldersTree *rbtree.RbTree) (map[string]interface{}, []string, map[string][]string) {
	var excludeFolders []string
	unexistFiles := make(map[string][]string)
	var includedFiles = make(map[string]interface{})

	foldersTree.Ascend(func(key *rbtree.Comparable) bool {
		info := (*key).(projectTreeNode).info
		if info.project == nil {
			return true
		}

		project := *info.projectPath

		// Add project base + exclude subfolder into exclude folders list
		parent := filepath.Dir(project)
		for _, s := range subfolderToExclude {
			sub := filepath.Join(parent, s)
			excludeFolders = append(excludeFolders, sub)
		}

		if info.project.OutputPaths != nil {
			for _, s := range info.project.OutputPaths {
				sub := filepath.Join(parent, s)
				excludeFolders = append(excludeFolders, sub)
			}
		}

		// Add compiles, contents and nones into included files map
		filesIncluded := getFilesIncludedIntoProject(info)
		for _, f := range filesIncluded {
			includedFiles[strings.ToUpper(f)] = nil
			if _, err := os.Stat(f); os.IsNotExist(err) {
				if found, ok := unexistFiles[project]; ok {
					found = append(found, f)
					unexistFiles[project] = found
				} else {
					unexistFiles[project] = []string{f}
				}
			}
		}
		return true
	})

	return includedFiles, excludeFolders, unexistFiles
}
