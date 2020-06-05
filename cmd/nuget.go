package cmd

import (
	"fmt"
	"github.com/aegoroff/godatastruct/rbtree"
	"solt/internal/msvc"
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
	Use:     "nu",
	Aliases: []string{"nuget"},
	Short:   "Get nuget packages information within projects or find Nuget mismatches in solution",
	RunE: func(cmd *cobra.Command, args []string) error {
		foldersTree := msvc.ReadSolutionDir(sourcesPath, appFileSystem)

		findNugetMismatches, err := cmd.Flags().GetBool(mismatchParamName)

		if err != nil {
			return err
		}

		if findNugetMismatches {
			showMismatches(foldersTree)
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

func showMismatches(foldersTree rbtree.RbTree) {

	solutions := msvc.SelectSolutions(foldersTree)

	var solutionProjects = make(map[string][]*msvc.FolderContent)

	// Each found solution
	allSolutionPaths := make(map[string]Matcher)
	for _, sln := range solutions {
		h := msvc.SelectAllSolutionProjectPaths(sln, normalize)
		allSolutionPaths[sln.Path] = NewExactMatchHS(&h)
		solutionProjects[sln.Path] = []*msvc.FolderContent{}
	}

	msvc.WalkProjects(foldersTree, func(prj *msvc.MsbuildProject, fold *msvc.Folder) {
		for k, v := range allSolutionPaths {
			if v.Match(normalize(prj.Path)) {
				solutionProjects[k] = append(solutionProjects[k], fold.Content)
			}
		}
	})

	mismatches := calculateMismatches(solutionProjects)

	if len(mismatches) == 0 {
		return
	}

	fmt.Println(" Different nuget package's versions in the same solution found:")

	const format = "  %v\t%v\n"
	tw := new(tabwriter.Writer).Init(appWriter, 0, 8, 4, ' ', 0)

	for sol, m := range mismatches {
		fmt.Printf("\n %s\n", sol)
		_, _ = fmt.Fprintf(tw, format, "Package", "Versions")
		_, _ = fmt.Fprintf(tw, format, "-------", "--------")
		for _, item := range m {
			_, _ = fmt.Fprintf(tw, format, item.pkg, strings.Join(item.versions, ", "))
		}
		_ = tw.Flush()
	}
}

func calculateMismatches(solutionProjects map[string][]*msvc.FolderContent) map[string][]*mismatch {
	var mismatches = make(map[string][]*mismatch)
	for sol, contents := range solutionProjects {
		var packagesMap = make(map[string][]string)
		for _, cnt := range contents {
			if cnt.Packages == nil && len(cnt.Projects) == 0 {
				continue
			}

			nugetPackages := getNugetPackages(cnt)

			for _, pkg := range nugetPackages {
				if v, ok := packagesMap[pkg.ID]; !ok {
					packagesMap[pkg.ID] = []string{pkg.Version}
				} else {
					// Only unique versions added
					if contains(v, pkg.Version) {
						continue
					}

					packagesMap[pkg.ID] = append(v, pkg.Version)
				}
			}
		}

		for pkg, vers := range packagesMap {
			// If one version it's OK (no mismatches)
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

func showPackagesInfoByFolders(foldersTree rbtree.RbTree) {
	const format = "  %v\t%v\n"
	tw := new(tabwriter.Writer).Init(appWriter, 0, 8, 4, ' ', 0)

	foldersTree.WalkInorder(func(n rbtree.Node) {
		folder := n.Key().(*msvc.Folder)
		content := folder.Content
		if content.Packages == nil && len(content.Projects) == 0 {
			return
		}

		nugetPackages := getNugetPackages(content)

		if len(nugetPackages) == 0 {
			return
		}

		parent := folder.Path
		fmt.Printf(" %s\n", parent)
		_, _ = fmt.Fprintf(tw, format, "Package", "Version")
		_, _ = fmt.Fprintf(tw, format, "-------", "--------")

		for _, p := range nugetPackages {
			_, _ = fmt.Fprintf(tw, format, p.ID, p.Version)
		}

		_ = tw.Flush()
		fmt.Println()
	})
}

func getNugetPackages(content *msvc.FolderContent) []msvc.NugetPackage {
	var nugetPackages []msvc.NugetPackage
	if content.Packages != nil {
		for _, p := range content.Packages.Packages {
			n := msvc.NugetPackage{ID: p.ID, Version: p.Version}
			nugetPackages = append(nugetPackages, n)
		}
	}
	for _, prj := range content.Projects {
		if prj.Project.PackageReferences == nil {
			continue
		}

		for _, p := range prj.Project.PackageReferences {
			n := msvc.NugetPackage{ID: p.ID, Version: p.Version}
			nugetPackages = append(nugetPackages, n)
		}
	}

	return nugetPackages
}
