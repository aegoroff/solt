package info

import (
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/dustin/go-humanize"
	"strings"
)

type totals struct {
	solutions    int64
	projects     int64
	projectTypes rbtree.RbTree
}

type typeStat struct {
	name      string
	count     int64
	solutions int64
}

func newTypeStat(name string) *typeStat {
	return &typeStat{name: name}
}

func (t *typeStat) Less(y rbtree.Comparable) bool {
	return strings.Compare(t.name, y.(*typeStat).name) < 0
}

func (t *typeStat) Equal(y rbtree.Comparable) bool {
	return strings.Compare(t.name, y.(*typeStat).name) == 0
}

func (t *typeStat) Count() string {
	return humanize.Comma(t.count)
}

func (t *typeStat) Solutions() string {
	return humanize.Comma(t.solutions)
}
