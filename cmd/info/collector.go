package info

import (
	"github.com/aegoroff/godatastruct/rbtree"
	"solt/msvc"
)

type collector struct {
	result *totals
	grp    *projectGroupper
}

func newCollector(grp *projectGroupper) *collector {
	return &collector{
		result: &totals{
			projectTypes: rbtree.New(),
		},
		grp: grp,
	}
}

func (c *collector) Solution(*msvc.VisualStudioSolution) {
	c.result.solutions++
	for k, v := range c.groupped() {
		c.updateType(k, int64(v))
	}
}

func (t *collector) updateType(k string, v int64) {
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

func (t *collector) groupped() map[string]int {
	return t.grp.ByType()
}
