package cmd

import (
	"fmt"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/spf13/afero"
	"log"
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
	Use:     "nuget",
	Aliases: []string{"nu"},
	Short:   "Get nuget packages information within projects or find Nuget mismatches in solution",
	RunE: func(cmd *cobra.Command, args []string) error {
		var solutions []string
		foldersTree := readProjectDir(sourcesPath, appFileSystem, func(we *walkEntry) {
			ext := strings.ToLower(filepath.Ext(we.Name))
			if ext == solutionFileExt {
				solutions = append(solutions, filepath.Join(we.Parent, we.Name))
			}
		})

		findNugetMismatches, err := cmd.Flags().GetBool(mismatchParamName)

		if err != nil {
			return err
		}

		if findNugetMismatches {
			showMismatches(solutions, foldersTree, appFileSystem)
		} else {
			showPackagesInfoByFolders(foldersTree)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(nugetCmd)

	nugetCmd.Flags().BoolP(mismatchParamName, "m", false, "Find packages to consolidate i.e. packages with different versions in the same solution")
}

func showMismatches(solutions []string, foldersTree *rbtree.RbTree, fs afero.Fs) {

	solutionProjects := getProjectsOfSolutions(solutions, foldersTree, fs)

	mismatches := calculateMismatches(solutionProjects)

	if len(mismatches) == 0 {
		return
	}

	fmt.Println(" Different nuget package's versions in the same solution found:")

	const format = "  %v\t%v\n"
	tw := new(tabwriter.Writer).Init(appWriter, 0, 8, 4, ' ', 0)

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

func getProjectsOfSolutions(solutions []string, foldersTree *rbtree.RbTree, fs afero.Fs) map[string][]*folderInfo {
	var solutionProjects = make(map[string][]*folderInfo)
	for _, sol := range solutions {
		f, err := fs.Open(sol)
		if err != nil {
			log.Println(err)
			continue
		}

		sln, err := solution.Parse(f)

		if err != nil {
			closeResource(f)
			log.Println(err)
			continue
		}

		var solutionProjectIds = make(StringHashSet)
		for _, sp := range sln.Projects {
			solutionProjectIds.Add(sp.Id)
		}

		foldersTree.WalkInorder(func(n *rbtree.Node) {
			finfo := (*n.Key).(projectTreeNode).info
			if finfo.project == nil {
				return
			}

			if solutionProjectIds.Contains(finfo.project.Id) {
				if v, ok := solutionProjects[sol]; !ok {
					solutionProjects[sol] = []*folderInfo{finfo}
				} else {
					solutionProjects[sol] = append(v, finfo)
				}
			}
		})
		closeResource(f)
	}
	return solutionProjects
}

func calculateMismatches(solutionProjects map[string][]*folderInfo) map[string][]*mismatch {
	var mismatches = make(map[string][]*mismatch)
	for sol, projects := range solutionProjects {
		var packagesMap = make(map[string][]string)
		for _, prj := range projects {
			if prj.packages == nil && (prj.project == nil || prj.project.PackageReferences == nil) {
				continue
			}

			nugetPackages := getNugetPackages(prj)

			for _, pkg := range nugetPackages {
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
	tw := new(tabwriter.Writer).Init(appWriter, 0, 8, 4, ' ', 0)

	foldersTree.WalkInorder(func(n *rbtree.Node) {
		fi := (*n.Key).(projectTreeNode).info
		if fi.packages == nil && (fi.project == nil || fi.project.PackageReferences == nil) {
			return
		}

		nugetPackages := getNugetPackages(fi)

		parent := filepath.Dir(*fi.projectPath)
		fmt.Printf(" %s\n", parent)
		fmt.Fprintf(tw, format, "Package", "Version")
		fmt.Fprintf(tw, format, "-------", "--------")

		for _, p := range nugetPackages {
			fmt.Fprintf(tw, format, p.Id, p.Version)
		}

		tw.Flush()
		fmt.Println()
	})
}

func getNugetPackages(fi *folderInfo) []nugetPackage {
	var nugetPackages []nugetPackage
	if fi.packages != nil {
		for _, p := range fi.packages.Packages {
			n := nugetPackage{Id: p.Id, Version: p.Version}
			nugetPackages = append(nugetPackages, n)
		}
	}
	if fi.project != nil && fi.project.PackageReferences != nil {
		for _, p := range fi.project.PackageReferences {
			n := nugetPackage{Id: p.Id, Version: p.Version}
			nugetPackages = append(nugetPackages, n)
		}
	}
	return nugetPackages
}
