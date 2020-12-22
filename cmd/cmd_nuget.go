package cmd

import (
	c9s "github.com/aegoroff/godatastruct/collections"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/spf13/cobra"
	"path/filepath"
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
		createCmd: func() executor {
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
	nugets := getFolderNugetPacks(foldersTree)

	prn := newNugetPrinter(p)
	it := rbtree.NewWalkInorder(nugets)

	it.Foreach(func(n rbtree.Comparable) {
		f := n.(*nugetFolder)
		prn.print(f.path, f.packs)
	})
}

func getFolderNugetPacks(foldersTree rbtree.RbTree) rbtree.RbTree {
	result := rbtree.NewRbTree()
	msvc.WalkProjectFolders(foldersTree, func(prj *msvc.MsbuildProject, fold *msvc.Folder) {
		packages := fold.Content.NugetPackages()

		var packs []*pack
		for _, np := range packages {
			versions := make(c9s.StringHashSet)
			versions.Add(np.Version)
			p := pack{
				pkg:      np.ID,
				versions: versions,
			}
			packs = append(packs, &p)
		}

		if len(packs) > 0 {
			n := newNugetFolder(fold.Path, packs)
			result.Insert(n)
		}
	})

	return result
}

func nugetBySolutions(foldersTree rbtree.RbTree, onlyMismatch bool, p printer) {
	solutions := msvc.SelectSolutions(foldersTree)

	// Each found solution
	allSolutionPaths := make(map[string][]string, len(solutions))
	for _, sln := range solutions {
		projects := sln.AllProjectPaths()
		allSolutionPaths[sln.Path] = getDirectories(projects)
	}

	nugets := getFolderNugetPacks(foldersTree)

	packs := getNugetPacks(allSolutionPaths, nugets)

	if onlyMismatch {
		keepOnlyMismatch(packs)
	}

	printNugetBySolutions(solutions, packs, onlyMismatch, p)
}

func getDirectories(paths []string) []string {
	result := make([]string, 0, len(paths))
	for _, path := range paths {
		result = append(result, filepath.Dir(path))
	}
	return result
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

func getNugetPacks(allSolPaths map[string][]string, nugets rbtree.RbTree) map[string][]*pack {
	var result = make(map[string][]*pack, len(allSolPaths))

	for spath, paths := range allSolPaths {
		for _, path := range paths {
			sv := newNugetFolder(path, nil)
			folder, ok := nugets.Search(sv)
			if ok {
				packs := folder.(*nugetFolder).packs
				result[spath] = append(result[spath], packs...)
			}
		}

		reduced := mergeNugetPacks(result[spath])
		if len(reduced) > 0 {
			result[spath] = reduced
		}
	}

	return result
}

func mergeNugetPacks(packs []*pack) []*pack {
	reduced := make([]*pack, 0, len(packs))
	unique := make(map[string]*pack)
	for _, p := range packs {
		exist, ok := unique[p.pkg]
		if ok {
			for _, v := range p.versions.Items() {
				exist.versions.Add(v)
			}
		} else {
			unique[p.pkg] = p
		}
	}

	for _, p := range unique {
		reduced = append(reduced, p)
	}
	return reduced
}

// keepOnlyMismatch removes all packs but only those
// which have more then one version on a nuget package
func keepOnlyMismatch(in map[string][]*pack) {
	toRemove := make(c9s.StringHashSet)
	for s, packs := range in {
		mismatchOnly := onlyMismatches(packs)
		if len(mismatchOnly) == 0 {
			toRemove.Add(s)
		} else {
			in[s] = mismatchOnly
		}
	}
	for s := range toRemove {
		delete(in, s)
	}
}

func onlyMismatches(packs []*pack) []*pack {
	filtered := make([]*pack, 0)
	for _, p := range packs {
		if p.versions.Count() > 1 {
			filtered = append(filtered, p)
		}
	}
	return filtered
}
