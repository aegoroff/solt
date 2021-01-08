package msvc

import (
	"github.com/aegoroff/godatastruct/rbtree"
)

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
