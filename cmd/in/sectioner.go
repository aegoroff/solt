package in

import (
	"solt/solution"
	"strings"

	c9s "github.com/aegoroff/godatastruct/collections"
)

type sectioner struct {
	configurations c9s.HashSet[string]
	platforms      c9s.HashSet[string]
}

func newSectioner() *sectioner {
	return &sectioner{
		configurations: c9s.NewHashSet[string](),
		platforms:      c9s.NewHashSet[string](),
	}
}

func (*sectioner) allow(section *solution.Section) bool {
	return section.Name == "SolutionConfigurationPlatforms"
}

func (c *sectioner) run(section *solution.Section) {
	for _, item := range section.Items {
		parts := strings.Split(item.Key, "|")
		c.configurations.Add(parts[0])
		c.platforms.Add(parts[1])
	}
}
