package msvc

import (
	"github.com/google/btree"
)

type walker interface {
	onFolder(f *Folder)
}

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

func walk(foldersTree *btree.BTree, walkers ...walker) {
	foldersTree.Ascend(func(n btree.Item) bool {
		fold := n.(*Folder)
		for _, w := range walkers {
			w.onFolder(fold)
		}
		return true
	})
}
