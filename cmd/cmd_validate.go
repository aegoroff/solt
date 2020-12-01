package cmd

import (
	"github.com/spf13/cobra"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"solt/msvc"
)

func newValidate() *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "va",
		Aliases: []string{"validate"},
		Short:   "Validates `SDK projects within solotion(s)",
		RunE: func(cmd *cobra.Command, args []string) error {
			foldersTree := msvc.ReadSolutionDir(sourcesPath, appFileSystem)

			solutions := msvc.SelectSolutions(foldersTree)

			for _, sol := range solutions {
				sln := sol.Solution

				g := simple.NewDirectedGraph()
				ids := make(map[string]graph.Node)
				ix := int64(1)
				for _, prj := range sln.Projects {
					n := newProjectNode(ix, prj.Path)
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
