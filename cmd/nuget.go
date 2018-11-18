package cmd

import (
	"fmt"
	"github.com/aegoroff/godatastruct/rbtree"
	"os"
	"path/filepath"
	"solt/solution"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

type mismatch struct {
	pkg      string
	versions []string
}

const mismatchParamName = "mismatch"

// nugetCmd represents the nuget command
var nugetCmd = &cobra.Command{
	Use:   "nuget",
	Short: "Get nuget packages information within projects or find Nuget mismatches in solution",
	Run: func(cmd *cobra.Command, args []string) {
		var solutions []string
		foldersTree := readProjectDir(sourcesPath, func(we *walkEntry) {
			ext := strings.ToLower(filepath.Ext(we.Name))
			if ext == solutionFileExt {
				solutions = append(solutions, filepath.Join(we.Parent, we.Name))
			}
		})

		findNugetMismatches, _ := cmd.Flags().GetBool(mismatchParamName)

		if findNugetMismatches {
			showMismatches(solutions, foldersTree)
		} else {
			showPackagesInfoByFolders(foldersTree)
		}

	},
}

func init() {
	rootCmd.AddCommand(nugetCmd)

	nugetCmd.Flags().BoolP(mismatchParamName, "m", false, "Find packages to consolidate i.e. packages with different versions in the same solution")
}

func showMismatches(solutions []string, foldersTree *rbtree.RbTree) {

	solutionProjects := getProjectsOfSolutions(solutions, foldersTree)

	mismatches := calculateMismatches(solutionProjects)

	if len(mismatches) == 0 {
		return
	}

	fmt.Println(" Different nuget package's versions in the same solution found:")

	const format = "  %v\t%v\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 4, ' ', 0)

	for sol, m := range mismatches {
		fmt.Printf("\n %s\n", sol)
		fmt.Fprintf(tw, format, "Package", "Versions")
		fmt.Fprintf(tw, format, "-------", "--------")
		for _, item := range m {
			fmt.Fprintf(tw, format, item.pkg, strings.Join(item.versions, ", "))
		}
		tw.Flush()
	}
}

func getProjectsOfSolutions(solutions []string, foldersTree *rbtree.RbTree) map[string][]*folderInfo {
	var solutionProjects = make(map[string][]*folderInfo)
	for _, sol := range solutions {
		sln, _ := solution.Parse(sol)
		var solutionProjectIds = make(map[string]interface{})
		for _, sp := range sln.Projects {
			solutionProjectIds[sp.Id] = nil
		}

		foldersTree.WalkInorder(func(n *rbtree.Node) {
			finfo := (*n.Key).(projectTreeNode).info
			if finfo.project == nil {
				return
			}

			if _, ok := solutionProjectIds[finfo.project.Id]; ok {
				if v, ok := solutionProjects[sol]; !ok {
					solutionProjects[sol] = []*folderInfo{finfo}
				} else {
					solutionProjects[sol] = append(v, finfo)
				}
			}
		})
	}
	return solutionProjects
}

func calculateMismatches(solutionProjects map[string][]*folderInfo) map[string][]*mismatch {
	var mismatches = make(map[string][]*mismatch)
	for sol, projects := range solutionProjects {
		var packagesMap = make(map[string][]string)
		for _, prj := range projects {
			if prj.packages == nil {
				continue
			}

			for _, pkg := range prj.packages.Packages {
				if v, ok := packagesMap[pkg.Id]; !ok {
					packagesMap[pkg.Id] = []string{pkg.Version}
				} else {
					if contains(v, pkg.Version) {
						continue
					}

					packagesMap[pkg.Id] = append(v, pkg.Version)
				}
			}
		}

		for pkg, vers := range packagesMap {
			if len(vers) < 2 {
				continue
			}

			m := mismatch{
				pkg:      pkg,
				versions: vers,
			}

			if v, ok := mismatches[sol]; !ok {
				mismatches[sol] = []*mismatch{&m}
			} else {
				mismatches[sol] = append(v, &m)
			}
		}
	}
	return mismatches
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func showPackagesInfoByFolders(foldersTree *rbtree.RbTree) {
	const format = "  %v\t%v\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 4, ' ', 0)

	foldersTree.WalkInorder(func(n *rbtree.Node) {
		v := (*n.Key).(projectTreeNode).info
		if v.packages == nil {
			return
		}

		parent := filepath.Dir(*v.projectPath)
		fmt.Printf(" %s\n", parent)
		fmt.Fprintf(tw, format, "Package", "Version")
		fmt.Fprintf(tw, format, "-------", "--------")

		for _, p := range v.packages.Packages {
			fmt.Fprintf(tw, format, p.Id, p.Version)
		}
		tw.Flush()
		fmt.Println()
	})
}
