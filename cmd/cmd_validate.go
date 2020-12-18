package cmd

import (
	"bytes"
	"encoding/xml"
	"fmt"
	c9s "github.com/aegoroff/godatastruct/collections"
	"github.com/akutz/sortfold"
	"github.com/spf13/cobra"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
	"io"
	"log"
	"path/filepath"
	"solt/internal/sys"
	"solt/msvc"
	"solt/solution"
	"unicode/utf8"
)

type validateCommand struct {
	baseCommand
	remove bool
}

type sdkProjectReference struct {
	Path string `xml:"Include,attr"`
}

func newValidate(c conf) *cobra.Command {
	var remove bool

	cc := cobraCreator{
		createCmd: func() command {
			return &validateCommand{
				baseCommand: newBaseCmd(c),
				remove:      remove,
			}
		},
	}

	cmd := cc.newCobraCommand("va", "validate", "Validates SDK projects within solution(s)")
	cmd.Flags().BoolVarP(&remove, "remove", "r", false, "Remove redundant project references from projects")

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

		if c.remove {
			c.updateProjects(refs)
		}
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

func (c *validateCommand) updateProjects(refs map[string]c9s.StringHashSet) {
	if len(refs) == 0 {
		return
	}

	filer := sys.NewFiler(c.fs, c.prn.writer())

	for project, rrs := range refs {
		ends := c.getElementsEnds(project, rrs)
		newContent := c.removeRedundantRefencesFromProject(project, ends)
		filer.Write(project, newContent)
	}
}

func (c *validateCommand) getElementsEnds(project string, toRemove c9s.StringHashSet) []int64 {
	f, err := c.fs.Open(filepath.Clean(project))
	defer sys.Close(f)
	if err != nil {
		return nil
	}

	decoder := xml.NewDecoder(f)
	pdir := filepath.Dir(project)

	ends := make([]int64, 0)
	for {
		t, err := decoder.Token()
		if t == nil {
			if err != nil && err != io.EOF {
				fmt.Println(err)
			}
			break
		}

		off := decoder.InputOffset()

		switch v := t.(type) {
		case xml.StartElement:
			// If we just read a StartElement token
			if v.Name.Local == "ProjectReference" {
				var prj sdkProjectReference
				// decode a whole chunk of following XML into the variable
				err = decoder.DecodeElement(&prj, &v)
				if err != nil {
					log.Println(err)
				} else {
					referenceFullPath := filepath.Join(pdir, prj.Path)
					if toRemove.Contains(referenceFullPath) {
						ends = append(ends, off)
					}
				}
			}
		}
	}

	return ends
}

func (c *validateCommand) removeRedundantRefencesFromProject(project string, ends []int64) []byte {
	f, err := c.fs.Open(filepath.Clean(project))
	defer sys.Close(f)
	if err != nil {
		return nil
	}

	buf := bytes.NewBuffer(nil)
	written, err := io.Copy(buf, f)

	if err != nil {
		return nil
	}

	result := make([]byte, 0, written)
	start := 0
	for _, end := range ends {
		n := int(end) - start
		start = int(end)

		portion := buf.Next(n)
		l := len(portion)
		for i := len(portion) - 2; i >= 0; i-- {
			r, _ := utf8.DecodeRune([]byte{portion[i]})
			if r == '>' || r == '\n' {
				l = i
				break
			}
		}

		portion = portion[:l]
		result = append(result, portion...)
	}

	result = append(result, buf.Next(buf.Len())...)

	return result
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
