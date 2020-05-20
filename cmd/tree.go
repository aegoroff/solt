package cmd

import (
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/akutz/sortfold"
	"strings"
)

type folderContent struct {
	packages  *Packages
	projects  []*msbuildProject
	solutions []*visualStudioSolution
}

type folder struct {
	content *folderContent
	path    string
}

func (x *folder) LessThan(y interface{}) bool {
	return sortfold.CompareFold(x.path, (y.(*folder)).path) < 0
}

func (x *folder) EqualTo(y interface{}) bool {
	return strings.EqualFold(x.path, (y.(*folder)).path)
}

func newTreeNode(f *folder) *rbtree.Comparable {
	var r rbtree.Comparable
	r = f
	return &r
}

func walkProjects(foldersTree *rbtree.RbTree, action func(prj *msbuildProject, fold *folder)) {
	foldersTree.WalkInorder(func(n *rbtree.Node) {
		fold := (*n.Key).(*folder)
		content := fold.content
		if len(content.projects) == 0 {
			return
		}

		// All found projects
		for _, prj := range content.projects {
			action(prj, fold)
		}
	})
}

func selectSolutions(foldersTree *rbtree.RbTree) []*visualStudioSolution {
	var solutions []*visualStudioSolution
	// Select only folders that contain solution(s)
	foldersTree.WalkInorder(func(n *rbtree.Node) {
		f := (*n.Key).(*folder)
		content := f.content
		if len(content.solutions) == 0 {
			return
		}
		for _, sln := range content.solutions {
			solutions = append(solutions, sln)
		}
	})
	return solutions
}
