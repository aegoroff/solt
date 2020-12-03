package cmd

import (
	"github.com/spf13/cobra"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
	"path/filepath"
	"solt/msvc"
	"solt/solution"
)

func newValidate() *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "va",
		Aliases: []string{"validate"},
		Short:   "Validates `SDK projects within solution(s)",
		RunE: func(cmd *cobra.Command, args []string) error {
			foldersTree := msvc.ReadSolutionDir(sourcesPath, appFileSystem)

			solutions, allProjects := msvc.SelectSolutionsAndProjects(foldersTree)

			prjMap := make(map[string]*msvc.MsbuildProject)

			for _, project := range allProjects {
				if !project.Project.IsSdkProject() {
					continue
				}
				prjMap[normalize(project.Path)] = project
			}

			for _, sol := range solutions {
				sln := sol.Solution
				solutionPath := filepath.Dir(sol.Path)

				g := simple.NewDirectedGraph()
				nodes := make(map[string]*projectNode)
				ix := int64(1)
				for _, prj := range sln.Projects {
					if prj.TypeID == solution.IDSolutionFolder {
						continue
					}

					fullProjectPath := normalize(filepath.Join(solutionPath, prj.Path))

					var msbuild *msvc.MsbuildProject
					msbuild, ok := prjMap[fullProjectPath]
					if !ok {
						continue
					}

					n := newProjectNode(ix, msbuild)
					nodes[fullProjectPath] = n
					ix++
					g.AddNode(n)
				}

				var roots []*projectNode
				for _, to := range nodes {
					if to.project.Project.ProjectReferences == nil {
						roots = append(roots, to)
					} else {
						dir := filepath.Dir(to.project.Path)

						for _, pref := range to.project.Project.ProjectReferences {
							full := filepath.Join(dir, pref.Path)
							from, ok := nodes[normalize(full)]
							if ok {
								e := g.NewEdge(from, to)
								g.SetEdge(e)
							}
						}
					}
				}
				ap := path.DijkstraAllPaths(g)
				for _, root := range roots {
					for _, to := range nodes {
						paths, _ := ap.AllBetween(root.ID(), to.ID())
						if paths != nil && len(paths) > 1 && root.ID() != to.ID() {
							appPrinter.cprint("from: %s to %s\n", root, to)
							for _, p := range paths {
								appPrinter.cprint("  %s\n", p)
							}
						}
					}
				}
				////cycles := topo.DirectedCyclesIn(g)
				//sorted, _ := topo.Sort(g)
				//for _, n := range sorted {
				//	dominators := flow.DominatorsSLT(n, g)
				//	dominators.Root()
				//}
			}

			return nil
		},
	}
	return cmd
}
