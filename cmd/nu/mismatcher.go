package nu

import (
	c9s "github.com/aegoroff/godatastruct/collections"
	"github.com/aegoroff/godatastruct/rbtree"
)

type mismatcher struct {
	nugets  rbtree.RbTree
	counter c9s.HashSet[string]
}

func newMismatcher(nugets rbtree.RbTree) *mismatcher {
	return &mismatcher{nugets: nugets, counter: c9s.NewHashSet[string]()}
}

func (m *mismatcher) count() int64 {
	return int64(m.counter.Count())
}

func (m *mismatcher) mismatchedPacks(mismatches []*pack, allPaths []string) rbtree.RbTree {
	result := rbtree.New()

	for _, mismatch := range mismatches {
		packs := make([]*pack, 0)
		for _, path := range allPaths {
			packs = append(packs, m.filter(path, mismatch)...)
		}
		if mismatch.versions.Count() > 1 {
			m.counter.Add(mismatch.pkg)
			node := newNugetFolder(mismatch.pkg, packs, nil)
			result.Insert(node)
		}
	}

	return result
}

func (m *mismatcher) filter(folderToSearch string, mismatch *pack) []*pack {
	searchKey := newNugetFolder(folderToSearch, nil, nil)
	found, ok := m.nugets.Search(searchKey)
	packs := make([]*pack, 0)

	if ok {
		nf := found.(*nugetFolder)
		for _, p := range nf.packs {
			if mismatch.match(p) {
				packs = append(packs, newPack(folderToSearch, p.versions.Items()...))
			}
		}
	}

	return packs
}
