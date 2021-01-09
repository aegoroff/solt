package nuget

import (
	"fmt"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/spf13/cobra"
	"path/filepath"
	"solt/cmd/fw"
	"solt/msvc"
	"strings"
)

const empiricNugetPacksForEachProject = 16

type nugetCommand struct {
	*fw.BaseCommand
	mismatch bool
	verbose  bool
}

type nugetByProjectCommand struct {
	*fw.BaseCommand
}

// New creates new command that does nuget packages feature
func New(c *fw.Conf) *cobra.Command {
	var mismatch bool
	var verbose bool

	cc := fw.NewCobraCreator(c, func() fw.Executor {
		exe := &nugetCommand{
			BaseCommand: fw.NewBaseCmd(c),
			mismatch:    mismatch,
			verbose:     verbose,
		}
		return fw.NewExecutorShowHelp(exe, c)
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

func newNugetByProject(c *fw.Conf) *cobra.Command {
	cc := fw.NewCobraCreator(c, func() fw.Executor {
		exe := &nugetByProjectCommand{fw.NewBaseCmd(c)}
		return fw.NewExecutorShowHelp(exe, c)
	})

	msg := "Get nuget packages information by projects' folders i.e. from packages.config or SDK project files"
	cmd := cc.NewCommand("p", "project", msg)

	return cmd
}

func (c *nugetCommand) Execute(*cobra.Command) error {
	foldersTree := msvc.ReadSolutionDir(c.SourcesPath(), c.Fs())
	c.execute(foldersTree)
	return nil
}

func (c *nugetByProjectCommand) Execute(*cobra.Command) error {
	foldersTree := msvc.ReadSolutionDir(c.SourcesPath(), c.Fs())
	c.execute(foldersTree)
	return nil
}

func newNugetFoldersTree(foldersTree rbtree.RbTree) rbtree.RbTree {
	result := rbtree.New()
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

	pSolution := newNugetPrinter(c.Prn(), c, "Package", 2)

	it := rbtree.NewAscend(packs)
	m := newMismatcher(nugets)

	it.Foreach(func(n rbtree.Comparable) {
		f := n.(*folder)

		pSolution.print(f.path, f.packs)

		if c.verbose {
			pPack := newNugetPrinter(c.Prn(), c, "Project", 5)
			mtree := m.mismatchedPacks(f.packs, f.sources)
			pPack.printTree(mtree, func(nf *folder) string {
				return fmt.Sprintf("Package: %s", nf.path)
			})
		}
	})
}

func (c *nugetByProjectCommand) execute(foldersTree rbtree.RbTree) {
	nugets := newNugetFoldersTree(foldersTree)

	prn := newNugetPrinter(c.Prn(), c, "Package", 2)

	prn.printTree(nugets, func(nf *folder) string {
		src := strings.Join(nf.sources, ", ")
		return fmt.Sprintf("%s (%s)", nf.path, src)
	})
}

// spreadNugetPacks binds all found nuget packages by solutions
func spreadNugetPacks(solutions []*msvc.VisualStudioSolution, nugets rbtree.RbTree) rbtree.RbTree {
	result := rbtree.New()

	for _, sol := range solutions {
		npacks, projectFolders := onlySolutionPacks(sol, nugets)
		reduced := mergeNugetPacks(npacks)

		if len(reduced) > 0 {
			nf := newNugetFolder(sol.Path(), reduced, projectFolders)
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
		val, ok := nugets.Search(sv)
		if ok {
			packs := val.(*folder).packs
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
	empty := make([]*folder, 0)

	rbtree.NewWalkInorder(in).Foreach(func(n rbtree.Comparable) {
		nf := n.(*folder)
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
	n := 0
	for _, p := range packs {
		if p.versions.Count() > 1 {
			packs[n] = p
			n++
		}
	}
	packs = packs[:n]
	return packs
}
