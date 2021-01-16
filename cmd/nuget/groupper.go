package nuget

import (
	"github.com/aegoroff/godatastruct/rbtree"
	"path/filepath"
	"solt/msvc"
)

type groupper struct {
	nugets   rbtree.RbTree
	groupped rbtree.RbTree
}

func newGroupper(nugets rbtree.RbTree) *groupper {
	return &groupper{nugets: nugets, groupped: rbtree.New()}
}

func (g *groupper) result(onlyMismatch bool) rbtree.RbTree {
	if onlyMismatch {
		g.keepOnlyMismatch()
	}
	return g.groupped
}

func (g *groupper) Solution(vs *msvc.VisualStudioSolution) {
	npacks, projectFolders := g.onlySolutionPacks(vs)
	reduced := mergeNugetPacks(npacks)

	if len(reduced) > 0 {
		nf := newNugetFolder(vs.Path(), reduced, projectFolders)
		g.groupped.Insert(nf)
	}
}

func (g *groupper) onlySolutionPacks(sol *msvc.VisualStudioSolution) ([]*pack, []string) {
	paths := sol.AllProjectPaths(filepath.Dir)
	npacks := make([]*pack, 0, len(paths)*empiricNugetPacksForEachProject)
	projectFolders := make([]string, 0, len(paths))

	for _, path := range paths {
		sv := newNugetFolder(path, nil, nil)
		val, ok := g.nugets.Search(sv)
		if ok {
			packs := val.(*folder).packs
			npacks = append(npacks, packs...)
			projectFolders = append(projectFolders, path)
		}
	}
	return npacks, projectFolders
}

func mergeNugetPacks(packs []*pack) []*pack {
	unique := make(map[string]*pack)
	for _, p := range packs {
		exist, ok := unique[p.pkg]
		if ok {
			exist.versions.AddRange(p.versions.Items()...)
		} else {
			unique[p.pkg] = p.copy()
		}
	}

	reduced := make([]*pack, len(unique))
	i := 0
	for _, p := range unique {
		reduced[i] = p
		i++
	}
	return reduced
}

// keepOnlyMismatch removes all packs but only those
// which have more then one version on a nuget package
func (g *groupper) keepOnlyMismatch() {
	empty := make([]*folder, 0)

	rbtree.NewWalkInorder(g.groupped).Foreach(func(n rbtree.Comparable) {
		nf := n.(*folder)
		mismatchOnly := onlyMismatches(nf.packs)
		if len(mismatchOnly) == 0 {
			empty = append(empty, nf)
		} else {
			nf.packs = mismatchOnly
		}
	})

	for _, n := range empty {
		g.groupped.Delete(n)
	}
}

func onlyMismatches(packs []*pack) []*pack {
	n := 0
	for _, p := range packs {
		if p.versions.Count() > 1 {
			packs[n] = p
			n++
		}
	}
	packs = packs[:n]
	return packs
}
