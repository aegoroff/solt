package cmd

import (
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/spf13/cobra"
	"path/filepath"
	"solt/msvc"
)

const empiricNugetPacksForEachProject = 16

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
	cmd.Flags().BoolVarP(&byProject, "project", "r", false, "Show packages by projects' folders instead")

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

func nugetBySolutions(foldersTree rbtree.RbTree, onlyMismatch bool, p printer) {
	nugets := getFolderNugetPacks(foldersTree)

	solutions := msvc.SelectSolutions(foldersTree)

	packs := getNugetPacks(solutions, nugets)

	if onlyMismatch {
		keepOnlyMismatch(packs)
	}

	printNugetBySolutions(packs, onlyMismatch, p)
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
			p := newPack(np.ID, np.Version)
			packs = append(packs, p)
		}

		if len(packs) > 0 {
			n := newNugetFolder(fold.Path, packs)
			result.Insert(n)
		}
	})

	return result
}

func printNugetBySolutions(packs rbtree.RbTree, onlyMismatch bool, p printer) {
	if packs.Len() == 0 {
		return
	}

	if onlyMismatch {
		p.cprint(" <red>Different nuget package's versions in the same solution found:</>")
	}

	prn := newNugetPrinter(p)

	rbtree.NewWalkInorder(packs).Foreach(func(n rbtree.Comparable) {
		nf := n.(*nugetFolder)
		prn.print(nf.path, nf.packs)
	})
}

func getNugetPacks(solutions []*msvc.VisualStudioSolution, nugets rbtree.RbTree) rbtree.RbTree {
	result := rbtree.NewRbTree()

	for _, sol := range solutions {
		paths := getDirectories(sol.AllProjectPaths())
		spacks := make([]*pack, 0, len(paths)*empiricNugetPacksForEachProject)

		for _, path := range paths {
			sv := newNugetFolder(path, nil)
			folder, ok := nugets.Search(sv)
			if ok {
				packs := folder.(*nugetFolder).packs
				spacks = append(spacks, packs...)
			}
		}

		reduced := mergeNugetPacks(spacks)

		if len(reduced) > 0 {
			nf := newNugetFolder(sol.Path, reduced)
			result.Insert(nf)
		}
	}

	return result
}

func getDirectories(paths []string) []string {
	result := make([]string, 0, len(paths))
	for _, path := range paths {
		result = append(result, filepath.Dir(path))
	}
	return result
}

func mergeNugetPacks(packs []*pack) []*pack {
	reduced := make([]*pack, 0, len(packs))
	unique := make(map[string]*pack)
	for _, p := range packs {
		exist, ok := unique[p.pkg]
		if ok {
			exist.versions.AddRange(p.versions.Items()...)
		} else {
			unique[p.pkg] = copyPack(p)
		}
	}

	for _, p := range unique {
		reduced = append(reduced, p)
	}
	return reduced
}

// keepOnlyMismatch removes all packs but only those
// which have more then one version on a nuget package
func keepOnlyMismatch(in rbtree.RbTree) {
	empty := make([]*nugetFolder, 0)

	rbtree.NewWalkInorder(in).Foreach(func(n rbtree.Comparable) {
		nf := n.(*nugetFolder)
		mismatchOnly := onlyMismatches(nf.packs)
		if len(mismatchOnly) == 0 {
			empty = append(empty, nf)
		} else {
			nf.packs = mismatchOnly
		}
	})

	for _, n := range empty {
		in.DeleteNode(n)
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
