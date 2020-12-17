package cmd

import (
	c9s "github.com/aegoroff/godatastruct/collections"
	"github.com/akutz/sortfold"
	"github.com/spf13/cobra"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
	"path/filepath"
	"solt/msvc"
	"solt/solution"
)

type validateCommand struct {
	baseCommand
}

func newValidate(c conf) *cobra.Command {
	cc := cobraCreator{
		createCmd: func() command {
			vac := validateCommand{
				baseCommand: newBaseCmd(c),
			}
			return &vac
		},
	}

	cmd := cc.newCobraCommand("va", "validate", "Validates SDK projects within solution(s)")

	return cmd
}

func (c *validateCommand) execute() error {
	foldersTree := msvc.ReadSolutionDir(c.sourcesPath, c.fs)

	solutions, allProjects := msvc.SelectSolutionsAndProjects(foldersTree)

	prjMap := newSdkProjects(allProjects)

	for _, sol := range solutions {
		g, nodes := newSolutionGraph(sol, prjMap)

		refs := findRedundantProjectReferences(g, nodes)
		printRedundantRefs(sol.Path, refs, c.prn)
	}

	return nil
}

func newSdkProjects(allProjects []*msvc.MsbuildProject) map[string]*msvc.MsbuildProject {
	prjMap := make(map[string]*msvc.MsbuildProject)

	for _, project := range allProjects {
		if !project.Project.IsSdkProject() {
			continue
		}
		prjMap[normalize(project.Path)] = project
	}
	return prjMap
}

func newSolutionGraph(sln *msvc.VisualStudioSolution, prjMap map[string]*msvc.MsbuildProject) (*simple.DirectedGraph, map[string]*projectNode) {
	solutionPath := filepath.Dir(sln.Path)
	g := simple.NewDirectedGraph()
	nodes := make(map[string]*projectNode)
	ix := int64(1)
	for _, prj := range sln.Solution.Projects {
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
		refs := getReferences(to, nodes)
		for _, ref := range refs {
			e := g.NewEdge(ref, to)
			g.SetEdge(e)
		}
	}
	return g, nodes
}

func findRedundantProjectReferences(g *simple.DirectedGraph, nodes map[string]*projectNode) map[string]c9s.StringHashSet {
	allPaths := path.DijkstraAllPaths(g)
	result := make(map[string]c9s.StringHashSet)
	for _, project := range nodes {
		refs := getReferences(project, nodes)

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
			result[project.String()] = rrs
		}
	}

	return result
}

func printRedundantRefs(solutionPath string, refs map[string]c9s.StringHashSet, p printer) {
	if len(refs) == 0 {
		return
	}
	p.cprint(" Solution: <green>%s</>\n", solutionPath)

	projects := make([]string, 0, len(refs))
	for s := range refs {
		projects = append(projects, s)
	}

	sortfold.Strings(projects)

	for _, project := range projects {
		p.cprint("   project: <bold>%s</> has redundant references\n", project)
		rrs := refs[project]
		items := rrs.Items()
		sortfold.Strings(items)
		for _, s := range items {
			p.cprint("     <gray>%s</>\n", s)
		}
	}
}

func getReferences(to *projectNode, nodes map[string]*projectNode) []*projectNode {
	if to.project.Project.ProjectReferences == nil {
		return []*projectNode{}
	}

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
