package cmd

import (
	c9s "github.com/aegoroff/godatastruct/collections"
	"github.com/spf13/afero"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
	"path/filepath"
	"solt/msvc"
	"solt/solution"
)

type sdkProjectsModule struct {
	prn         printer
	fs          afero.Fs
	sourcesPath string
	h           sdkModuleHandler
}

func newSdkProjectsModule(fs afero.Fs, p printer, sourcesPath string, h sdkModuleHandler) *sdkProjectsModule {
	return &sdkProjectsModule{
		prn:         p,
		fs:          fs,
		sourcesPath: sourcesPath,
		h:           h,
	}
}

func (m *sdkProjectsModule) execute() {
	foldersTree := msvc.ReadSolutionDir(m.sourcesPath, m.fs)

	solutions, allProjects := msvc.SelectSolutionsAndProjects(foldersTree)

	sdkProjects := m.onlySdkProjects(allProjects)

	for _, sol := range solutions {
		g, nodes := m.newSolutionGraph(sol, sdkProjects)

		refs := m.redundantRefs(g, nodes)

		m.h.onRedundantRefs(sol.Path, refs)
	}
}

func (*sdkProjectsModule) onlySdkProjects(allProjects []*msvc.MsbuildProject) map[string]*msvc.MsbuildProject {
	prjMap := make(map[string]*msvc.MsbuildProject)

	for _, project := range allProjects {
		if !project.Project.IsSdkProject() {
			continue
		}
		prjMap[normalize(project.Path)] = project
	}
	return prjMap
}

func (m *sdkProjectsModule) newSolutionGraph(sln *msvc.VisualStudioSolution, prjMap map[string]*msvc.MsbuildProject) (*simple.DirectedGraph, map[string]*projectNode) {
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
		refs := m.getReferences(to, nodes)
		for _, ref := range refs {
			e := g.NewEdge(ref, to)
			g.SetEdge(e)
		}
	}
	return g, nodes
}

func (m *sdkProjectsModule) redundantRefs(g *simple.DirectedGraph, nodes map[string]*projectNode) map[string]c9s.StringHashSet {
	allPaths := path.DijkstraAllPaths(g)
	result := make(map[string]c9s.StringHashSet)
	for _, project := range nodes {
		refs := m.getReferences(project, nodes)

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

func (*sdkProjectsModule) getReferences(to *projectNode, nodes map[string]*projectNode) []*projectNode {
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
