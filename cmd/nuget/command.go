package nuget

import (
	"fmt"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/spf13/cobra"
	"path/filepath"
	"solt/cmd/api"
	"solt/msvc"
	"strings"
)

const empiricNugetPacksForEachProject = 16

type nugetCommand struct {
	api.BaseCommand
	mismatch bool
	verbose  bool
}

type nugetByProjectCommand struct {
	api.BaseCommand
}

// New creates new command that does nuget packages feature
func New(c *api.Conf) *cobra.Command {
	var mismatch bool
	var verbose bool

	cc := api.NewCobraCreator(c, func() api.Executor {
		return &nugetCommand{
			BaseCommand: api.NewBaseCmd(c),
			mismatch:    mismatch,
			verbose:     verbose,
		}
	})

	descr := "Get nuget packages information within solutions"
	cmd := cc.NewCommand("nu", "nuget", descr)

	mdescr := "Find packages to consolidate i.e. packages with different versions in the same solution"
	cmd.Flags().BoolVarP(&mismatch, "mismatch", "m", false, mdescr)

	vdescr := "Output details about mismatched nuget packages"
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, vdescr)

	cmd.AddCommand(newNugetByProject(c))

	return cmd
}

func newNugetByProject(c *api.Conf) *cobra.Command {
	cc := api.NewCobraCreator(c, func() api.Executor {
		return &nugetByProjectCommand{
			BaseCommand: api.NewBaseCmd(c),
		}
	})

	msg := "Get nuget packages information by projects' folders i.e. from packages.config or SDK project files"
	cmd := cc.NewCommand("p", "project", msg)

	return cmd
}

func (c *nugetCommand) Execute() error {
	foldersTree := msvc.ReadSolutionDir(c.SourcesPath(), c.Fs())
	c.execute(foldersTree)
	return nil
}

func (c *nugetByProjectCommand) Execute() error {
	foldersTree := msvc.ReadSolutionDir(c.SourcesPath(), c.Fs())
	c.execute(foldersTree)
	return nil
}

func newNugetFoldersTree(foldersTree rbtree.RbTree) rbtree.RbTree {
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

func (c *nugetCommand) execute(foldersTree rbtree.RbTree) {
	nugets := newNugetFoldersTree(foldersTree)

	solutions := msvc.SelectSolutions(foldersTree)

	packs := spreadNugetPacks(solutions, nugets)

	if c.mismatch {
		keepOnlyMismatch(packs)
	}

	if packs.Len() == 0 {
		return
	}

	if c.mismatch {
		c.Prn().Cprint(" <red>Different nuget package's versions in the same solution found:</>\n")
	}

	pSolution := newNugetPrinter(c.Prn(), "Package", 2)

	it := rbtree.NewAscend(packs)
	m := newMismatcher(nugets)

	it.Foreach(func(n rbtree.Comparable) {
		f := n.(*nugetFolder)

		pSolution.print(f.path, f.packs)

		if c.verbose {
			pPack := newNugetPrinter(c.Prn(), "Project", 5)
			mtree := m.mismatchedPacks(f.packs, f.sources)
			pPack.printTree(mtree, func(nf *nugetFolder) string {
				return fmt.Sprintf("Package: %s", nf.path)
			})
		}
	})
}

func (c *nugetByProjectCommand) execute(foldersTree rbtree.RbTree) {
	nugets := newNugetFoldersTree(foldersTree)

	prn := newNugetPrinter(c.Prn(), "Package", 2)

	prn.printTree(nugets, func(nf *nugetFolder) string {
		src := strings.Join(nf.sources, ", ")
		return fmt.Sprintf("%s (%s)", nf.path, src)
	})
}

// spreadNugetPacks binds all found nuget packages by solutions
func spreadNugetPacks(solutions []*msvc.VisualStudioSolution, nugets rbtree.RbTree) rbtree.RbTree {
	result := rbtree.NewRbTree()

	for _, sol := range solutions {
		npacks, projectFolders := onlySolutionPacks(sol, nugets)
		reduced := mergeNugetPacks(npacks)

		if len(reduced) > 0 {
			nf := newNugetFolder(sol.Path, reduced, projectFolders)
			result.Insert(nf)
		}
	}

	return result
}

func onlySolutionPacks(sol *msvc.VisualStudioSolution, nugets rbtree.RbTree) ([]*pack, []string) {
	paths := sol.AllProjectPaths(filepath.Dir)
	npacks := make([]*pack, 0, len(paths)*empiricNugetPacksForEachProject)
	projectFolders := make([]string, 0, len(paths))

	for _, path := range paths {
		sv := newNugetFolder(path, nil, nil)
		folder, ok := nugets.Search(sv)
		if ok {
			packs := folder.(*nugetFolder).packs
			npacks = append(npacks, packs...)
			projectFolders = append(projectFolders, path)
		}
	}
	return npacks, projectFolders
}

func mergeNugetPacks(packs []*pack) []*pack {
	unique := make(map[string]*pack)
	for _, p := range packs {
		exist, ok := unique[p.pkg]
		if ok {
			exist.versions.AddRange(p.versions.Items()...)
		} else {
			unique[p.pkg] = p.copy()
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
