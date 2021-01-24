package va

import (
	"github.com/aegoroff/godatastruct/rbtree"
	"solt/msvc"
)

type sdkProjects struct {
	tree rbtree.RbTree
}

func newSdkProjects(allProjects []*msvc.MsbuildProject) *sdkProjects {
	tree := rbtree.New()

	for _, p := range allProjects {
		if p.IsSdkProject() {
			tree.Insert(p)
		}
	}

	return &sdkProjects{
		tree: tree,
	}
}

func (s *sdkProjects) search(p *msvc.MsbuildProject) (*msvc.MsbuildProject, bool) {
	found, ok := s.tree.Search(p)
	if ok {
		return found.(*msvc.MsbuildProject), true
	}
	return nil, false
}
