package cmd

import (
	"github.com/spf13/cobra"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"path/filepath"
	"solt/msvc"
	"solt/solution"
)

func newValidate() *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "va",
		Aliases: []string{"validate"},
		Short:   "Validates `SDK projects within solotion(s)",
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
				ids := make(map[string]graph.Node)
				ix := int64(1)
				for _, prj := range sln.Projects {
					if prj.TypeID == solution.IDSolutionFolder {
						continue
					}

					fullProjectPath := filepath.Join(solutionPath, prj.Path)

					var msbuild *msvc.MsbuildProject
					msbuild, ok := prjMap[normalize(fullProjectPath)]
					if !ok {
						continue
					}

					n := newProjectNode(ix, msbuild)
					ids[prj.Path] = n
					ix++
					g.AddNode(n)
				}
			}

			return nil
		},
	}
	return cmd
}
