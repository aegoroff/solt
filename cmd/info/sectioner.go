package info

import (
	c9s "github.com/aegoroff/godatastruct/collections"
	"solt/solution"
	"strings"
)

type configurationPlatform struct {
	configurations c9s.StringHashSet
	platforms      c9s.StringHashSet
}

func newConfigurationPlatform() *configurationPlatform {
	return &configurationPlatform{
		configurations: make(c9s.StringHashSet),
		platforms:      make(c9s.StringHashSet),
	}
}

func (c *configurationPlatform) allow(section *solution.Section) bool {
	return section.Name == "SolutionConfigurationPlatforms"
}

func (c *configurationPlatform) run(section *solution.Section) {
	for _, item := range section.Items {
		parts := strings.Split(item.Key, "|")
		c.configurations.Add(parts[0])
		c.platforms.Add(parts[1])
	}
}
