package info

import (
	c9s "github.com/aegoroff/godatastruct/collections"
	"github.com/spf13/cobra"
	"solt/cmd/api"
	"solt/msvc"
	"solt/solution"
	"strconv"
	"strings"
)

type infoCommand struct {
	*api.BaseCommand
	margin int
}

// New creates new command that shows information about solutions
func New(c *api.Conf) *cobra.Command {
	cc := api.NewCobraCreator(c, func() api.Executor {
		return &infoCommand{
			BaseCommand: api.NewBaseCmd(c),
			margin:      2,
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

		tbl := api.NewTabler(c.Prn(), c.margin)

		tbl.AddLine("Header", sln.Header)
		tbl.AddLine("Product", sln.Comment)
		tbl.AddLine("Visual Studio Version", sln.VisualStudioVersion)
		tbl.AddLine("Minimum Visual Studio Version", sln.MinimumVisualStudioVersion)

		tbl.Print()

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

	tbl := api.NewTabler(c.Prn(), c.margin)
	tbl.AddHead("Project type", "Count")

	for k, v := range byType {
		tbl.AddLine(k, strconv.Itoa(v))
	}
	tbl.Print()
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

	prn := newPrinter(c.margin, p)

	prn.print(configurations, "Configuration")

	p.Cprint("\n")

	prn.print(platforms, "Platform")

	p.Cprint("\n")
}
