package in

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

func (c *collector) updateType(k string, v int64) {
	c.result.projects += v
	key := newTypeStat(k)
	n, ok := c.result.projectTypes.Search(newTypeStat(k))
	if ok {
		ts := n.(*typeStat)
		ts.solutions++
		ts.count += v
	} else {
		key.solutions = 1
		key.count = v
		c.result.projectTypes.Insert(key)
	}
}

func (c *collector) groupped() map[string]int {
	return c.grp.ByType()
}
