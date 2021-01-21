package info

import (
	"github.com/aegoroff/godatastruct/rbtree"
	"solt/msvc"
)

type totaler struct {
	result *totals
	grp    *projectGroupper
}

func newTotaler(grp *projectGroupper) *totaler {
	return &totaler{
		result: &totals{
			projectTypes: rbtree.New(),
		},
		grp: grp,
	}
}

func (t *totaler) Solution(*msvc.VisualStudioSolution) {
	t.result.solutions++
	for k, v := range t.groupped() {
		t.updateType(k, int64(v))
	}
}

func (t *totaler) updateType(k string, v int64) {
	t.result.projects += v
	key := newTypeStat(k)
	n, ok := t.result.projectTypes.Search(newTypeStat(k))
	if ok {
		ts := n.(*typeStat)
		ts.solutions++
		ts.count += v
	} else {
		key.solutions = 1
		key.count = v
		t.result.projectTypes.Insert(key)
	}
}

func (t *totaler) groupped() map[string]int {
	return t.grp.ByType()
}
