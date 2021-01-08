package msvc

import (
	"github.com/aegoroff/godatastruct/rbtree"
)

// SelectProjects gets all Visual Studion solutions found in a directory
func SelectProjects(foldersTree rbtree.RbTree) []*MsbuildProject {
	var projects []*MsbuildProject

	WalkProjectFolders(foldersTree, func(project *MsbuildProject, folder *Folder) {
		projects = append(projects, project)
	})

	return projects
}

// SelectSolutions gets all Visual Studion solutions found in a directory
func SelectSolutions(foldersTree rbtree.RbTree) []*VisualStudioSolution {
	w := walkSol{solutions: make([]*VisualStudioSolution, 0)}
	walk(foldersTree, &w)
	return w.solutions
}

// SelectSolutionsAndProjects gets all Visual Studion solutions and projects found in a directory
func SelectSolutionsAndProjects(foldersTree rbtree.RbTree) ([]*VisualStudioSolution, []*MsbuildProject) {
	ws := walkSol{solutions: make([]*VisualStudioSolution, 0)}
	var projects []*MsbuildProject

	wp := walkPrj{handler: func(project *MsbuildProject, folder *Folder) {
		projects = append(projects, project)
	}}

	walk(foldersTree, &ws, &wp)
	return ws.solutions, projects
}

// WalkProjectFolders traverse all projects found in solution(s) folder
func WalkProjectFolders(foldersTree rbtree.RbTree, action ProjectHandler) {
	w := &walkPrj{handler: action}
	walk(foldersTree, w)
}

type walkPrj struct {
	handler ProjectHandler
}

func (w *walkPrj) walk(f *Folder) {
	content := f.Content

	if len(content.Projects) == 0 {
		return
	}

	// All found projects
	for _, prj := range content.Projects {
		w.handler(prj, f)
	}
}

type walkSol struct {
	solutions []*VisualStudioSolution
}

func (w *walkSol) walk(f *Folder) {
	content := f.Content
	// Select only folders that contain solution(s)
	if len(content.Solutions) == 0 {
		return
	}

	w.solutions = append(w.solutions, content.Solutions...)
}

func walk(foldersTree rbtree.RbTree, walkers ...walker) {
	rbtree.NewWalkInorder(foldersTree).Foreach(func(n rbtree.Comparable) {
		fold := n.(*Folder)
		for _, w := range walkers {
			w.walk(fold)
		}
	})
}
