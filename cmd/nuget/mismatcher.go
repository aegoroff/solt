package nuget

import "github.com/aegoroff/godatastruct/rbtree"

type mismatcher struct {
	nugets rbtree.RbTree
}

func newMismatcher(nugets rbtree.RbTree) *mismatcher {
	return &mismatcher{nugets: nugets}
}

func (m *mismatcher) mismatchedPacks(mismatches []*pack, allPaths []string) rbtree.RbTree {
	result := rbtree.NewRbTree()

	for _, mismatch := range mismatches {
		packs := make([]*pack, 0)
		for _, path := range allPaths {
			packs = append(packs, m.filter(path, mismatch)...)
		}
		node := newNugetFolder(mismatch.pkg, packs, nil)
		result.Insert(node)
	}

	return result
}

func (m *mismatcher) filter(folderToSearch string, mismatch *pack) []*pack {
	searchKey := newNugetFolder(folderToSearch, nil, nil)
	found, ok := m.nugets.Search(searchKey)
	packs := make([]*pack, 0)

	if ok {
		nf := found.(*folder)
		for _, p := range nf.packs {
			if mismatch.match(p) {
				packs = append(packs, newPack(folderToSearch, p.versions.Items()...))
			}
		}
	}

	return packs
}
