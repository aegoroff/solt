package info

import (
	"fmt"
	c9s "github.com/aegoroff/godatastruct/collections"
	"github.com/akutz/sortfold"
	"github.com/cheynewallace/tabby"
	"github.com/spf13/cobra"
	"solt/cmd/api"
	"solt/msvc"
	"solt/solution"
	"strings"
)

type infoCommand struct {
	api.BaseCommand
}

// New creates new command that shows information about solutions
func New(c *api.Conf) *cobra.Command {
	cc := api.NewCobraCreator(c, func() api.Executor {
		ic := infoCommand{
			api.NewBaseCmd(c),
		}
		return &ic
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

		t.AddLine("  Header", sln.Header)
		t.AddLine("  Product", sln.Comment)
		t.AddLine("  Visual Studio Version", sln.VisualStudioVersion)
		t.AddLine("  Minimum Visual Studio Version", sln.MinimumVisualStudioVersion)

		t.Print()

		fmt.Println()

		showProjectsInfo(sln.Projects, c.Prn())
		showSectionsInfo(sln.GlobalSections, c.Prn())
	}

	return nil
}

func showProjectsInfo(projects []*solution.Project, p api.Printer) {
	var byType = make(map[string]int)

	for _, p := range projects {
		byType[p.Type]++
	}

	const format = "  %v\t%v\n"

	p.Tprint(format, "Project type", "Count")
	p.Tprint(format, "------------", "-----")

	for k, v := range byType {
		p.Tprint(format, k, v)
	}
	p.Flush()
	fmt.Println()
}

func showSectionsInfo(sections []*solution.Section, p api.Printer) {
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

	const format = "  %v\n"

	p.Tprint(format, "Configuration")
	p.Tprint(format, "------------")

	sortedConfigurations := configurations.Items()
	sortfold.Strings(sortedConfigurations)

	for _, k := range sortedConfigurations {
		p.Tprint(format, k)
	}
	p.Flush()
	fmt.Println()

	p.Tprint(format, "Platform")
	p.Tprint(format, "--------")

	sortedPlatforms := platforms.Items()
	sortfold.Strings(sortedPlatforms)

	for _, k := range sortedPlatforms {
		p.Tprint(format, k)
	}
	p.Flush()
	p.Cprint("\n")
}
