package cmd

import (
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/spf13/cobra"
	"solt/msvc"
)

type nugetCommand struct {
	baseCommand
	mismatch  bool
	byProject bool
}

func newNuget(c *conf) *cobra.Command {
	var mismatch bool
	var byProject bool

	cc := cobraCreator{
		createCmd: func() command {
			nc := nugetCommand{
				baseCommand: newBaseCmd(c),
				mismatch:    mismatch,
				byProject:   byProject,
			}
			return &nc
		},
	}

	cmd := cc.newCobraCommand("nu", "nuget", "Get nuget packages information within solutions, projects or find Nuget mismatches in solution")

	cmd.Flags().BoolVarP(&mismatch, "mismatch", "m", false, "Find packages to consolidate i.e. packages with different versions in the same solution")
	cmd.Flags().BoolVarP(&byProject, "project", "r", false, "Show packages by projects instead")

	return cmd
}

func (c *nugetCommand) execute() error {
	foldersTree := msvc.ReadSolutionDir(c.sourcesPath, c.fs)

	if c.mismatch || !c.byProject {
		nugetBySolutions(foldersTree, c.mismatch, c.prn)
	} else {
		nugetByProjects(foldersTree, c.prn)
	}

	return nil
}

func nugetByProjects(foldersTree rbtree.RbTree, p printer) {
	prn := newNugetPrinter(p)
	msvc.WalkProjectFolders(foldersTree, func(prj *msvc.MsbuildProject, fold *msvc.Folder) {
		content := fold.Content
		pchan := make(chan *msvc.NugetPackage, 4)
		go readNugetPackages(content, pchan)

		var packs []*pack
		for np := range pchan {
			p := pack{
				pkg:      np.ID,
				versions: []string{np.Version},
			}
			packs = append(packs, &p)
		}

		if len(packs) > 0 {
			prn.print(fold.Path, packs)
		}
	})
}

func nugetBySolutions(foldersTree rbtree.RbTree, onlyMismatch bool, p printer) {
	solutions := msvc.SelectSolutions(foldersTree)

	var allProjectFolders = make(map[string]*msvc.FolderContent, foldersTree.Len())

	// Each found solution
	allSolutionPaths := make(map[string]Matcher, len(solutions))
	for _, sln := range solutions {
		projects := sln.AllProjectPaths()
		allSolutionPaths[sln.Path] = NewExactMatch(projects)
	}

	msvc.WalkProjectFolders(foldersTree, func(prj *msvc.MsbuildProject, fold *msvc.Folder) {
		allProjectFolders[prj.Path] = fold.Content
	})

	packs := getNugetPacks(allSolutionPaths, allProjectFolders, onlyMismatch)

	printNugetBySolutions(solutions, packs, onlyMismatch, p)
}

func printNugetBySolutions(solutions []*msvc.VisualStudioSolution, packs map[string][]*pack, onlyMismatch bool, p printer) {
	if len(packs) == 0 {
		return
	}

	if onlyMismatch {
		p.cprint(" <red>Different nuget package's versions in the same solution found:</>")
	}

	prn := newNugetPrinter(p)
	for _, sln := range solutions {
		if pks, ok := packs[sln.Path]; ok {
			prn.print(sln.Path, pks)
		}
	}
}

func getNugetPacks(allSolPaths map[string]Matcher, allPrjFolders map[string]*msvc.FolderContent, onlyMismatch bool) map[string][]*pack {
	allPkg := mapAllPackages(allPrjFolders)

	var result = make(map[string][]*pack, len(allSolPaths))

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

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func mapAllPackages(allPrjFolders map[string]*msvc.FolderContent) map[string]map[string]string {
	var packagesByProject = make(map[string]map[string]string, len(allPrjFolders))

	for ppath, content := range allPrjFolders {
		if len(content.Projects) == 0 {
			continue
		}

		var packagesMap map[string]string
		packagesMap = make(map[string]string)
		packagesByProject[ppath] = packagesMap

		pchan := make(chan *msvc.NugetPackage, 4)

		go readNugetPackages(content, pchan)

		for pkg := range pchan {
			packagesMap[pkg.ID] = pkg.Version
		}
	}
	return packagesByProject
}

func readNugetPackages(content *msvc.FolderContent, pchan chan<- *msvc.NugetPackage) {
	if content.Packages != nil {
		// old style projects (nuget packages references in separate files)
		for _, p := range content.Packages.Packages {
			pchan <- &msvc.NugetPackage{ID: p.ID, Version: p.Version}
		}
	}
	for _, prj := range content.Projects {
		if prj.Project.PackageReferences == nil {
			continue
		}

		// If SDK project nuget packages included into project file
		for _, p := range prj.Project.PackageReferences {
			pchan <- &msvc.NugetPackage{ID: p.ID, Version: p.Version}
		}
	}
	close(pchan)
}
