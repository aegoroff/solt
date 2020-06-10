package cmd

import (
	"fmt"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/akutz/sortfold"
	"github.com/spf13/cobra"
	"solt/msvc"
	"sort"
)

type nugetCmd struct {
	mismatch  bool
	byProject bool
}

func newNuget() *cobra.Command {
	opts := nugetCmd{}
	var cmd = &cobra.Command{
		Use:     "nu",
		Aliases: []string{"nuget"},
		Short:   "Get nuget packages information within solutions, projects or find Nuget mismatches in solution",
		RunE: func(cmd *cobra.Command, args []string) error {
			foldersTree := msvc.ReadSolutionDir(sourcesPath, appFileSystem)

			if opts.mismatch || !opts.byProject {
				nugetBySolutions(foldersTree, opts.mismatch)
			} else {
				nugetByProjects(foldersTree)
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&opts.mismatch, "mismatch", "m", false, "Find packages to consolidate i.e. packages with different versions in the same solution")
	cmd.Flags().BoolVarP(&opts.byProject, "project", "r", false, "Show packages by projects instead")

	return cmd
}

func nugetByProjects(foldersTree rbtree.RbTree) {
	prn := newNugetPrinter(appWriter)
	msvc.WalkProjectFolders(foldersTree, func(prj *msvc.MsbuildProject, fold *msvc.Folder) {
		content := fold.Content
		nugetPackages := getNugetPackages(content)

		if len(nugetPackages) == 0 {
			return
		}

		var packs []*pack
		for _, np := range nugetPackages {
			p := pack{
				pkg:      np.ID,
				versions: []string{np.Version},
			}
			packs = append(packs, &p)
		}

		prn.print(fold.Path, packs)
	})
}

func nugetBySolutions(foldersTree rbtree.RbTree, onlyMismatch bool) {
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

	packs := getNugetPacks(allSolutionPaths, allProjectFolders, onlyMismatch)

	printNugetBySolutions(solutions, packs, onlyMismatch)
}

func printNugetBySolutions(solutions []*msvc.VisualStudioSolution, packs map[string][]*pack, onlyMismatch bool) {
	if len(packs) == 0 {
		return
	}

	if onlyMismatch {
		_, _ = fmt.Fprintln(appWriter, " Different nuget package's versions in the same solution found:")
	}

	sort.Slice(solutions, func(i, j int) bool {
		return sortfold.CompareFold(solutions[i].Path, solutions[j].Path) < 0
	})

	prn := newNugetPrinter(appWriter)
	for _, sln := range solutions {
		if pks, ok := packs[sln.Path]; ok {
			prn.print(sln.Path, pks)
		}
	}
}

func getNugetPacks(allSolPaths map[string]Matcher, allPrjFolders map[string]*msvc.FolderContent, onlyMismatch bool) map[string][]*pack {
	allPkg := mapAllPackages(allPrjFolders)

	var result = make(map[string][]*pack)

	// Reduce packages
	for spath, match := range allSolPaths {
		packagesVers := mapPackagesInSolution(allPkg, match)

		// Reduce packages in solution
		packs := reducePacks(packagesVers, onlyMismatch)
		if len(packs) > 0 {
			result[spath] = packs
		}
	}

	return result
}

func reducePacks(packagesVers map[string][]string, onlyMismatch bool) []*pack {
	var result []*pack
	for pkg, vers := range packagesVers {
		// If one version it's OK (no mismatches)
		if onlyMismatch && len(vers) < 2 {
			continue
		}

		m := pack{
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

func getNugetPackages(content *msvc.FolderContent) []*msvc.NugetPackage {
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
