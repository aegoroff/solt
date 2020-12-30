package info

import (
	c9s "github.com/aegoroff/godatastruct/collections"
	"github.com/cheynewallace/tabby"
	"github.com/spf13/cobra"
	"solt/cmd/api"
	"solt/msvc"
	"solt/solution"
	"strings"
)

type infoCommand struct {
	*api.BaseCommand
	m *api.Marginer
}

// New creates new command that shows information about solutions
func New(c *api.Conf) *cobra.Command {
	cc := api.NewCobraCreator(c, func() api.Executor {
		return &infoCommand{
			BaseCommand: api.NewBaseCmd(c),
			m:           api.NewMarginer(2),
		}
	})

	cmd := cc.NewCommand("in", "info", "Get information about found solutions")
	return cmd
}

func (c *infoCommand) Execute() error {
	foldersTree := msvc.ReadSolutionDir(c.SourcesPath(), c.Fs())

	solutions := msvc.SelectSolutions(foldersTree)

	for _, sol := range solutions {
		sln := sol.Solution

		c.Prn().Cprint(" <gray>%s</>\n", sol.Path)

		t := tabby.NewCustom(c.Prn().Twriter())

		t.AddLine(c.m.Margin("Header"), sln.Header)
		t.AddLine(c.m.Margin("Product"), sln.Comment)
		t.AddLine(c.m.Margin("Visual Studio Version"), sln.VisualStudioVersion)
		t.AddLine(c.m.Margin("Minimum Visual Studio Version"), sln.MinimumVisualStudioVersion)

		t.Print()

		c.Prn().Cprint("\n")

		c.showProjectsInfo(sln.Projects, c.Prn())
		c.showSectionsInfo(sln.GlobalSections, c.Prn())
	}

	return nil
}

func (c *infoCommand) showProjectsInfo(projects []*solution.Project, p api.Printer) {
	var byType = make(map[string]int)

	for _, p := range projects {
		byType[p.Type]++
	}

	format := c.m.Margin("%v\t%v\n")

	p.Tprint(format, "Project type", "Count")
	p.Tprint(format, "------------", "-----")

	for k, v := range byType {
		p.Tprint(format, k, v)
	}
	p.Flush()
	p.Cprint("\n")
}

func (c *infoCommand) showSectionsInfo(sections []*solution.Section, p api.Printer) {
	var configurations = make(c9s.StringHashSet)
	var platforms = make(c9s.StringHashSet)

	for _, s := range sections {
		if s.Name != "SolutionConfigurationPlatforms" {
			continue
		}
		for _, item := range s.Items {
			parts := strings.Split(item.Key, "|")
			configuration := parts[0]
			platform := parts[1]
			configurations.Add(configuration)
			platforms.Add(platform)
		}
	}

	prn := newPrinter(c.m, p)

	prn.print(configurations, "Configuration")

	p.Cprint("\n")

	prn.print(platforms, "Platform")

	p.Cprint("\n")
}
