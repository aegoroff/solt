package cmd

import (
	"fmt"
	"github.com/aegoroff/godatastruct/rbtree"
	"solt/internal/msvc"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

type mismatch struct {
	pkg      string
	versions []string
}

type nugetPackages []*msvc.NugetPackage
type mismatches []*mismatch

func (n nugetPackages) Len() int           { return len(n) }
func (n nugetPackages) Less(i, j int) bool { return n[i].ID < n[j].ID }
func (n nugetPackages) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }

func (m mismatches) Len() int           { return len(m) }
func (m mismatches) Less(i, j int) bool { return m[i].pkg < m[j].pkg }
func (m mismatches) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }

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

	var allProjectFolders = make(map[string]*msvc.FolderContent)

	// Each found solution
	allSolutionPaths := make(map[string]Matcher)
	for _, sln := range solutions {
		h := msvc.SelectAllSolutionProjectPaths(sln, normalize)
		allSolutionPaths[sln.Path] = NewExactMatchHS(&h)
	}

	msvc.WalkProjectFolders(foldersTree, func(prj *msvc.MsbuildProject, fold *msvc.Folder) {
		allProjectFolders[normalize(prj.Path)] = fold.Content
	})

	mismatches := calculateMismatches(allSolutionPaths, allProjectFolders)

	if len(mismatches) == 0 {
		return
	}

	_, _ = fmt.Fprintln(appWriter, " Different nuget package's versions in the same solution found:")

	const format = "  %v\t%v\n"
	tw := new(tabwriter.Writer).Init(appWriter, 0, 8, 4, ' ', 0)

	for sol, m := range mismatches {
		_, _ = fmt.Fprintf(appWriter, "\n %s\n", sol)
		_, _ = fmt.Fprintf(tw, format, "Package", "Versions")
		_, _ = fmt.Fprintf(tw, format, "-------", "--------")
		sort.Sort(m)
		for _, item := range m {
			_, _ = fmt.Fprintf(tw, format, item.pkg, strings.Join(item.versions, ", "))
		}
		_ = tw.Flush()
	}
}

func calculateMismatches(allSolPaths map[string]Matcher, allPrjFolders map[string]*msvc.FolderContent) map[string]mismatches {
	allPkg := mapAllPackages(allPrjFolders)

	var mismatches = make(map[string]mismatches)

	// Reduce packages
	for spath, match := range allSolPaths {
		packagesVers := mapPackagesInSolution(allPkg, match)

		// Reduce packages in solution
		mm := reducePackages(packagesVers)
		if len(mm) > 0 {
			mismatches[spath] = mm
		}
	}

	return mismatches
}

func reducePackages(packagesVers map[string][]string) mismatches {
	var result []*mismatch
	for pkg, vers := range packagesVers {
		// If one version it's OK (no mismatches)
		if len(vers) < 2 {
			continue
		}

		m := mismatch{
			pkg:      pkg,
			versions: vers,
		}

		result = append(result, &m)
	}
	return result
}

func mapPackagesInSolution(packagesByProject map[string]map[string]string, match Matcher) map[string][]string {
	packagesVers := make(map[string][]string)

	for ppath, pkg := range packagesByProject {
		if !match.Match(ppath) {
			continue
		}

		for pkg, ver := range pkg {
			if v, ok := packagesVers[pkg]; !ok {
				packagesVers[pkg] = []string{ver}
			} else {
				// Only unique versions added
				if contains(v, ver) {
					continue
				}

				packagesVers[pkg] = append(v, ver)
			}
		}
	}

	return packagesVers
}

func mapAllPackages(allPrjFolders map[string]*msvc.FolderContent) map[string]map[string]string {
	var packagesByProject = make(map[string]map[string]string)

	for ppath, content := range allPrjFolders {
		if len(content.Projects) == 0 {
			continue
		}

		var packagesMap map[string]string
		packagesMap = make(map[string]string)
		packagesByProject[ppath] = packagesMap

		nugetPackages := getNugetPackages(content)

		for _, pkg := range nugetPackages {
			packagesMap[pkg.ID] = pkg.Version
		}
	}
	return packagesByProject
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

	msvc.WalkProjectFolders(foldersTree, func(prj *msvc.MsbuildProject, fold *msvc.Folder) {
		content := fold.Content
		nugetPackages := getNugetPackages(content)

		if len(nugetPackages) == 0 {
			return
		}

		parent := fold.Path
		_, _ = fmt.Fprintf(appWriter, " %s\n", parent)
		_, _ = fmt.Fprintf(tw, format, "Package", "Version")
		_, _ = fmt.Fprintf(tw, format, "-------", "--------")

		sort.Sort(nugetPackages)
		for _, p := range nugetPackages {
			_, _ = fmt.Fprintf(tw, format, p.ID, p.Version)
		}

		_ = tw.Flush()
		_, _ = fmt.Fprintln(appWriter)
	})
}

func getNugetPackages(content *msvc.FolderContent) nugetPackages {
	var nugetPackages []*msvc.NugetPackage
	if content.Packages != nil {
		for _, p := range content.Packages.Packages {
			n := msvc.NugetPackage{ID: p.ID, Version: p.Version}
			nugetPackages = append(nugetPackages, &n)
		}
	}
	for _, prj := range content.Projects {
		if prj.Project.PackageReferences == nil {
			continue
		}

		for _, p := range prj.Project.PackageReferences {
			n := msvc.NugetPackage{ID: p.ID, Version: p.Version}
			nugetPackages = append(nugetPackages, &n)
		}
	}

	return nugetPackages
}
