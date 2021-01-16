package nuget

import (
	"fmt"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/spf13/cobra"
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
	nugets := newNugetFoldersTree(foldersTree)

	solutions := msvc.SelectSolutions(foldersTree)

	grp := newGroupper(nugets)
	fw.SolutionSlice(solutions).Foreach(grp)
	packs := grp.result(c.mismatch)

	if packs.Len() == 0 {
		return nil
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
	return nil
}

func (c *nugetByProjectCommand) Execute(*cobra.Command) error {
	foldersTree := msvc.ReadSolutionDir(c.SourcesPath(), c.Fs())
	nugets := newNugetFoldersTree(foldersTree)

	prn := newNugetPrinter(c.Prn(), c, "Package", 2)

	prn.printTree(nugets, func(nf *folder) string {
		src := strings.Join(nf.sources, ", ")
		return fmt.Sprintf("%s (%s)", nf.path, src)
	})
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
