package nu

import (
	"github.com/aegoroff/godatastruct/rbtree"
	"path/filepath"
	"solt/msvc"
)

type grouper struct {
	nugets  rbtree.RbTree
	grouped rbtree.RbTree
}

func newGroupper(nugets rbtree.RbTree) *grouper {
	return &grouper{nugets: nugets, grouped: rbtree.New()}
}

func (g *grouper) result(onlyMismatch bool) rbtree.RbTree {
	if onlyMismatch {
		g.keepOnlyMismatch()
	}
	return g.grouped
}

func (g *grouper) Solution(vs *msvc.VisualStudioSolution) {
	npacks, projectFolders := g.onlySolutionPacks(vs)
	reduced := mergeNugetPacks(npacks)

	if len(reduced) > 0 {
		nf := newNugetFolder(vs.Path(), reduced, projectFolders)
		g.grouped.Insert(nf)
	}
}

func (g *grouper) onlySolutionPacks(sol *msvc.VisualStudioSolution) ([]*pack, []string) {
	paths := sol.AllProjectPaths(filepath.Dir)
	npacks := make([]*pack, 0, len(paths)*empiricNugetPacksForEachProject)
	projectFolders := make([]string, 0, len(paths))

	for _, path := range paths {
		sv := newNugetFolder(path, nil, nil)
		val, ok := g.nugets.Search(sv)
		if ok {
			packs := val.(*nugetFolder).packs
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
// which have more then one version on a nu package
func (g *grouper) keepOnlyMismatch() {
	empty := make([]*nugetFolder, 0)

	rbtree.NewWalkInorder(g.grouped).Foreach(func(n rbtree.Comparable) {
		nf := n.(*nugetFolder)
		mismatchOnly := onlyMismatches(nf.packs)
		if len(mismatchOnly) == 0 {
			empty = append(empty, nf)
		} else {
			nf.packs = mismatchOnly
		}
	})

	for _, n := range empty {
		g.grouped.Delete(n)
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
