package nuget

import (
	c9s "github.com/aegoroff/godatastruct/collections"
	"strings"
)

// pack defines nuget package descriptor
type pack struct {
	pkg      string
	versions c9s.StringHashSet
}

func newPack(id string, versions ...string) *pack {
	vs := c9s.NewStringHashSet()
	vs.AddRange(versions...)
	return &pack{
		pkg:      id,
		versions: vs,
	}
}

func (p *pack) copy() *pack {
	return newPack(p.pkg, p.versions.Items()...)
}

func (p *pack) match(other *pack) bool {
	idEqual := strings.EqualFold(p.pkg, other.pkg)

	if !idEqual {
		return false
	}

	for v := range p.versions {
		if other.versions.Contains(v) {
			return true
		}
	}

	return false
}
