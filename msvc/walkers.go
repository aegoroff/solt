package msvc

import (
	"github.com/aegoroff/godatastruct/rbtree"
)

type walkPrj struct {
	handler ProjectHandler
}

func (w *walkPrj) onFolder(f *Folder) {
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

func (w *walkSol) onFolder(f *Folder) {
	content := f.Content
	// Select only folders that contain solution(s)
	if len(content.Solutions) == 0 {
		return
	}
	for _, sln := range content.Solutions {
		w.solutions = append(w.solutions, sln)
	}
}

func walk(foldersTree rbtree.RbTree, walkers ...walker) {
	rbtree.NewAscend(foldersTree).Iterate(func(n rbtree.Node) bool {
		fold := n.Key().(*Folder)
		for _, w := range walkers {
			w.onFolder(fold)
		}
		return true
	})
}
