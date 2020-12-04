package cmd

import (
	c9s "github.com/aegoroff/godatastruct/collections"
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
		Short:   "Validates SDK projects within solution(s)",
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

			var currentSolution string
			for _, sol := range solutions {
				sln := sol.Solution
				solutionPath := filepath.Dir(sol.Path)
				currentSolution = sol.Path

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

				for _, to := range nodes {
					if to.project.Project.ProjectReferences != nil {
						refs := getReferences(to, nodes)
						for _, ref := range refs {
							e := g.NewEdge(ref, to)
							g.SetEdge(e)
						}
					}
				}

				allPaths := path.DijkstraAllPaths(g)
				for _, node := range nodes {
					refs := getReferences(node, nodes)

					rrs := make(c9s.StringHashSet)

					for _, from := range refs {
						for _, to := range refs {
							if from.ID() == to.ID() {
								continue
							}
							paths, _ := allPaths.AllBetween(from.ID(), to.ID())
							if paths != nil && len(paths) > 0 {
								rrs.Add(from.String())
							}
						}
					}

					if rrs.Count() > 0 {
						if currentSolution != "" {
							appPrinter.cprint(" Solution: <green>%s</>\n", currentSolution)
							currentSolution = ""
						}
						appPrinter.cprint("   project: <bold>%s</> has redundant references\n", node)
						for s := range rrs {
							appPrinter.cprint("    <gray>%s</>\n", s)
						}
					}
				}
			}

			return nil
		},
	}
	return cmd
}

func getReferences(to *projectNode, nodes map[string]*projectNode) []*projectNode {
	dir := filepath.Dir(to.project.Path)

	var result []*projectNode
	for _, pref := range to.project.Project.ProjectReferences {
		full := filepath.Join(dir, pref.Path)
		from, ok := nodes[normalize(full)]
		if ok {
			result = append(result, from)
		}
	}
	return result
}
