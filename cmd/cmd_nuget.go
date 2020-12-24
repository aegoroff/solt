package cmd

import (
	"fmt"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/spf13/cobra"
	"path/filepath"
	"solt/msvc"
	"strings"
)

const empiricNugetPacksForEachProject = 16

type nugetCommand struct {
	baseCommand
	mismatch bool
}

type nugetByProjectCommand struct {
	baseCommand
}

func newNuget(c *conf) *cobra.Command {
	var mismatch bool

	cc := cobraCreator{
		createCmd: func() executor {
			nc := nugetCommand{
				baseCommand: newBaseCmd(c),
				mismatch:    mismatch,
			}
			return &nc
		},
		c: c,
	}

	descr := "Get nuget packages information within solutions"
	cmd := cc.newCobraCommand("nu", "nuget", descr)

	mdescr := "Find packages to consolidate i.e. packages with different versions in the same solution"
	cmd.Flags().BoolVarP(&mismatch, "mismatch", "m", false, mdescr)

	cmd.AddCommand(newNugetByProject(c))

	return cmd
}

func newNugetByProject(c *conf) *cobra.Command {
	cc := cobraCreator{
		createCmd: func() executor {
			return &nugetByProjectCommand{
				baseCommand: newBaseCmd(c),
			}
		},
		c: c,
	}

	msg := "Get nuget packages information by projects' folders i.e. from packages.config or SDK project files"
	cmd := cc.newCobraCommand("p", "project", msg)

	return cmd
}

func (c *nugetCommand) execute() error {
	foldersTree := msvc.ReadSolutionDir(c.sourcesPath, c.fs)
	nugetBySolutions(foldersTree, c.mismatch, c.prn)
	return nil
}

func (c *nugetByProjectCommand) execute() error {
	foldersTree := msvc.ReadSolutionDir(c.sourcesPath, c.fs)
	nugetByProjects(foldersTree, c.prn)
	return nil
}

func nugetBySolutions(foldersTree rbtree.RbTree, onlyMismatch bool, p printer) {
	nugets := getFolderNugetPacks(foldersTree)

	solutions := msvc.SelectSolutions(foldersTree)

	packs := spreadNugetPacks(solutions, nugets)

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
		src := strings.Join(f.sources, ", ")
		prn.print(fmt.Sprintf("<bold>%s</> (%s)", f.path, src), f.packs)
	})
}

func getFolderNugetPacks(foldersTree rbtree.RbTree) rbtree.RbTree {
	result := rbtree.NewRbTree()
	msvc.WalkProjectFolders(foldersTree, func(prj *msvc.MsbuildProject, fold *msvc.Folder) {
		packages, sources := fold.Content.NugetPackages()
		if len(packages) == 0 {
			return
		}

		packs := make([]*pack, len(packages))
		for i, np := range packages {
			packs[i] = newPack(np.ID, np.Version)
		}

		n := newNugetFolder(fold.Path, packs, sources)
		result.Insert(n)
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

// spreadNugetPacks binds all found nuget packages by solutions
func spreadNugetPacks(solutions []*msvc.VisualStudioSolution, nugets rbtree.RbTree) rbtree.RbTree {
	result := rbtree.NewRbTree()

	for _, sol := range solutions {
		raw := onlySolutionPacks(sol, nugets)
		reduced := mergeNugetPacks(raw)

		if len(reduced) > 0 {
			nf := newNugetFolder(sol.Path, reduced, nil)
			result.Insert(nf)
		}
	}

	return result
}

func onlySolutionPacks(sol *msvc.VisualStudioSolution, nugets rbtree.RbTree) []*pack {
	paths := getDirectories(sol.AllProjectPaths())
	result := make([]*pack, 0, len(paths)*empiricNugetPacksForEachProject)

	for _, path := range paths {
		sv := newNugetFolder(path, nil, nil)
		folder, ok := nugets.Search(sv)
		if ok {
			packs := folder.(*nugetFolder).packs
			result = append(result, packs...)
		}
	}
	return result
}

func getDirectories(paths []string) []string {
	result := paths[:0]
	for _, path := range paths {
		result = append(result, filepath.Dir(path))
	}
	return result
}

func mergeNugetPacks(packs []*pack) []*pack {
	unique := make(map[string]*pack)
	for _, p := range packs {
		exist, ok := unique[p.pkg]
		if ok {
			exist.versions.AddRange(p.versions.Items()...)
		} else {
			unique[p.pkg] = copyPack(p)
		}
	}

	reduced := make([]*pack, len(unique))
	i := 0
	for _, p := range unique {
		reduced[i] = p
		i++
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
	filtered := packs[:0]
	for _, p := range packs {
		if p.versions.Count() > 1 {
			filtered = append(filtered, p)
		}
	}
	return filtered
}
