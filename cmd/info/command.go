package info

import (
	c9s "github.com/aegoroff/godatastruct/collections"
	"github.com/spf13/cobra"
	"solt/cmd/fw"
	"solt/msvc"
	"solt/solution"
	"strconv"
	"strings"
)

type infoCommand struct {
	*fw.BaseCommand
	margin int
}

// New creates new command that shows information about solutions
func New(c *fw.Conf) *cobra.Command {
	cc := fw.NewCobraCreator(c, func() fw.Executor {
		exe := &infoCommand{
			BaseCommand: fw.NewBaseCmd(c),
			margin:      2,
		}
		return fw.NewExecutorShowHelp(exe, c)
	})

	cmd := cc.NewCommand("in", "info", "Get information about found solutions")
	return cmd
}

func (c *infoCommand) Execute(*cobra.Command) error {
	foldersTree := msvc.ReadSolutionDir(c.SourcesPath(), c.Fs())

	solutions := msvc.SelectSolutions(foldersTree)
	msvc.SortSolutions(solutions)

	for _, sol := range solutions {
		sln := sol.Solution

		c.Prn().Cprint(" <gray>%s</>\n", sol.Path)

		tbl := fw.NewTabler(c, c.margin)

		tbl.AddLine("Header", sln.Header)
		tbl.AddLine("Product", sln.Comment)
		tbl.AddLine("Visual Studio Version", sln.VisualStudioVersion)
		tbl.AddLine("Minimum Visual Studio Version", sln.MinimumVisualStudioVersion)

		tbl.Print()

		c.Prn().Println()

		c.showProjectsInfo(sln.Projects)
		c.showSectionsInfo(sln.GlobalSections)
	}

	return nil
}

func (c *infoCommand) showProjectsInfo(projects []*solution.Project) {
	var byType = make(map[string]int)

	for _, p := range projects {
		byType[p.Type]++
	}

	tbl := fw.NewTabler(c, c.margin)
	tbl.AddHead("Project type", "Count")

	for k, v := range byType {
		tbl.AddLine(k, strconv.Itoa(v))
	}
	tbl.Print()
	c.Prn().Println()
}

func (c *infoCommand) showSectionsInfo(sections []*solution.Section) {
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

	prn := newPrinter(c.margin, c)

	prn.print(configurations, "Configuration")

	c.Prn().Println()

	prn.print(platforms, "Platform")

	c.Prn().Println()
}
