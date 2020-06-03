package msvc

import (
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/akutz/sortfold"
	"strings"
)

// FolderContent defines a filesystem folder information about
// it's MSVC content (solutions, projects, etc.)
type FolderContent struct {
	Packages  *packages
	Projects  []*MsbuildProject
	Solutions []*VisualStudioSolution
}

// Folder defines filesystem folder descriptor (path and content structure)
type Folder struct {
	Content *FolderContent
	Path    string
}

// LessThan implements rbtree.Comparable interface
func (x *Folder) LessThan(y interface{}) bool {
	return sortfold.CompareFold(x.Path, (y.(*Folder)).Path) < 0
}

// EqualTo implements rbtree.Comparable interface
func (x *Folder) EqualTo(y interface{}) bool {
	return strings.EqualFold(x.Path, (y.(*Folder)).Path)
}

// String implements rbtree.Comparable interface
func (x *Folder) String() string {
	return x.Path
}

// WalkProjects traverse all projects found in solution(s) folder
func WalkProjects(foldersTree rbtree.RbTree, action func(prj *MsbuildProject, fold *Folder)) {
	foldersTree.WalkInorder(func(n rbtree.Node) {
		fold := n.Key().(*Folder)
		content := fold.Content
		if len(content.Projects) == 0 {
			return
		}

		// All found projects
		for _, prj := range content.Projects {
			action(prj, fold)
		}
	})
}

// SelectSolutions gets all Visual Studion solutions found in a directory
func SelectSolutions(foldersTree rbtree.RbTree) []*VisualStudioSolution {
	var solutions []*VisualStudioSolution
	// Select only folders that contain solution(s)
	foldersTree.WalkInorder(func(n rbtree.Node) {
		f := n.Key().(*Folder)
		content := f.Content
		if len(content.Solutions) == 0 {
			return
		}
		for _, sln := range content.Solutions {
			solutions = append(solutions, sln)
		}
	})
	return solutions
}
