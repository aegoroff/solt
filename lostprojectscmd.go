package main

import (
	"fmt"
	"github.com/aegoroff/godatastruct/rbtree"
	"gonum.org/v1/gonum/graph/simple"
	"os"
	"path/filepath"
	"solt/solution"
	"strings"
)

type projectSolution struct {
	project  string
	solution string
}

func lostprojectscmd(opt options) error {

	var solutions []string
	tree := readProjectDir(opt.Path, func(we *walkEntry) {
		ext := strings.ToLower(filepath.Ext(we.Name))
		if ext == solutionFileExt {
			sp := filepath.Join(we.Parent, we.Name)
			solutions = append(solutions, sp)
		}
	})

	allProjectsWithinSolutions, solutionGraph := getAllSolutionsProjects(solutions)

	extendSolutionGraph(solutionGraph, tree)

	projectsOutsideSolution, filesInsideSolution := getOutsideProjectsAndFilesInsideSolution(tree, allProjectsWithinSolutions)

	projectsOutside, projectsOutsideSolutionWithFilesInside := separateOutsideProjects(projectsOutsideSolution, filesInsideSolution)

	sortAndOutputToStdout(projectsOutside)

	if len(projectsOutsideSolutionWithFilesInside) > 0 {
		fmt.Printf("\nThese projects not included into any solution but files in the projects' folders used in another projects within a solution:\n\n")
	}

	sortAndOutputToStdout(projectsOutsideSolutionWithFilesInside)

	unexistProjects := getUnexistProjects(allProjectsWithinSolutions)

	if len(unexistProjects) > 0 {
		fmt.Printf("\nThese projects included into a solution but not found in the file system:\n")
	}

	outputSortedMapToStdout(unexistProjects, "Solution")

	return nil
}

func extendSolutionGraph(solutionGraph *simple.UndirectedGraph, tree *rbtree.RbTree) {
	var fileNodesMap = make(map[string]int64)
	for _, parentNode := range solutionGraph.Nodes() {
		path := parentNode.(*node).name
		tn := createProjectTreeNode(path, nil)
		found, ok := rbtree.Search(tree.Root, tn)
		if ok {
			info := (*found.Key).(projectTreeNode).info
			if info.project == nil {
				continue
			}

			files := getFilesIncludedIntoProject(info)

			for _, f := range files {
				if nodeId, ok := fileNodesMap[f]; ok {
					edge := solutionGraph.NewEdge(parentNode, solutionGraph.Node(nodeId))
					solutionGraph.SetEdge(edge)
					continue
				}

				n := solutionGraph.NewNode()
				fileNode := node{
					nodeID: n.ID(),
					name:   f,
				}
				fileNodesMap[f] = n.ID()
				solutionGraph.AddNode(&fileNode)

				edge := solutionGraph.NewEdge(parentNode, &fileNode)
				solutionGraph.SetEdge(edge)
			}
		}
	}
}

func getUnexistProjects(allProjectsWithinSolutions map[string]*projectSolution) map[string][]string {
	var result = make(map[string][]string)
	for _, prj := range allProjectsWithinSolutions {
		if _, err := os.Stat(prj.project); !os.IsNotExist(err) {
			continue
		}

		if found, ok := result[prj.solution]; ok {
			found = append(found, prj.project)
			result[prj.solution] = found
		} else {
			result[prj.solution] = []string{prj.project}
		}
	}
	return result
}

func getOutsideProjectsAndFilesInsideSolution(tree *rbtree.RbTree, allProjectsWithinSolutions map[string]*projectSolution) ([]*folderInfo, map[string]interface{}) {

	var projectsOutsideSolution []*folderInfo
	var filesInsideSolution = make(map[string]interface{})

	rbtree.WalkInorder(tree.Root, func(n *rbtree.Node) {
		info := (*n.Key).(projectTreeNode).info
		if info.project == nil {
			return
		}

		id := strings.ToUpper(info.project.Id)

		_, ok := allProjectsWithinSolutions[id]
		if !ok {
			projectsOutsideSolution = append(projectsOutsideSolution, info)
		} else {
			filesIncluded := getFilesIncludedIntoProject(info)

			for _, f := range filesIncluded {
				filesInsideSolution[strings.ToUpper(f)] = nil
			}
		}
	})

	return projectsOutsideSolution, filesInsideSolution
}

func separateOutsideProjects(projectsOutsideSolution []*folderInfo, filesInsideSolution map[string]interface{}) ([]string, []string) {
	var projectsOutside []string
	var projectsOutsideSolutionWithFilesInside []string
	for _, info := range projectsOutsideSolution {
		projectFiles := getFilesIncludedIntoProject(info)

		var includedIntoOther = false
		for _, f := range projectFiles {
			pf := strings.ToUpper(f)
			if _, ok := filesInsideSolution[pf]; !ok {
				continue
			}

			dir := filepath.Dir(*info.projectPath)

			if strings.Contains(pf, strings.ToUpper(dir)) {
				includedIntoOther = true
				break
			}
		}

		if !includedIntoOther {
			projectsOutside = append(projectsOutside, *info.projectPath)
		} else {
			projectsOutsideSolutionWithFilesInside = append(projectsOutsideSolutionWithFilesInside, *info.projectPath)
		}
	}
	return projectsOutside, projectsOutsideSolutionWithFilesInside
}

func getAllSolutionsProjects(solutions []string) (map[string]*projectSolution, *simple.UndirectedGraph) {
	graph := simple.NewUndirectedGraph()

	var projectsInSolution = make(map[string]*projectSolution)
	for _, solpath := range solutions {
		sln, _ := solution.Parse(solpath)

		n := graph.NewNode()
		solNode := node{
			nodeID: n.ID(),
			name:   solpath,
		}

		graph.AddNode(&solNode)

		for _, p := range sln.Projects {
			// Skip solution folders
			if p.TypeId == "{2150E333-8FDC-42A3-9474-1A3956D46DE8}" {
				continue
			}

			id := strings.ToUpper(p.Id)

			// Already added
			if _, ok := projectsInSolution[id]; ok {
				continue
			}

			parent := filepath.Dir(solpath)
			pp := filepath.Join(parent, p.Path)

			n := graph.NewNode()
			prjNode := node{
				nodeID: n.ID(),
				name:   pp,
			}

			graph.AddNode(&prjNode)
			edge := graph.NewEdge(&solNode, &prjNode)
			graph.SetEdge(edge)

			projectsInSolution[id] = &projectSolution{
				project:  pp,
				solution: solpath,
			}
		}
	}
	return projectsInSolution, graph
}
