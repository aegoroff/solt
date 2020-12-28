package cmd

import (
	"fmt"
	c9s "github.com/aegoroff/godatastruct/collections"
	"github.com/akutz/sortfold"
	"github.com/cheynewallace/tabby"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"solt/cmd/api"
	"solt/msvc"
	"solt/solution"
	"strings"
)

type infoCommand struct {
	baseCommand
}

func newInfo(c *conf) *cobra.Command {
	cc := cobraCreator{
		createCmd: func() api.Executor {
			ic := infoCommand{
				newBaseCmd(c),
			}
			return &ic
		},
		c: c,
	}

	cmd := cc.NewCobraCommand("in", "info", "Get information about found solutions")
	return cmd
}

func (c *infoCommand) Execute() error {
	foldersTree := msvc.ReadSolutionDir(c.sourcesPath, c.fs)

	solutions := msvc.SelectSolutions(foldersTree)

	for _, sol := range solutions {
		sln := sol.Solution

		c.prn.SetColor(color.FgGray)
		c.prn.Cprint(" %s\n", sol.Path)
		c.prn.ResetColor()

		t := tabby.NewCustom(c.prn.Twriter())

		t.AddLine("  Header", sln.Header)
		t.AddLine("  Product", sln.Comment)
		t.AddLine("  Visual Studio Version", sln.VisualStudioVersion)
		t.AddLine("  Minimum Visual Studio Version", sln.MinimumVisualStudioVersion)

		t.Print()

		fmt.Println()

		showProjectsInfo(sln.Projects, c.prn)
		showSectionsInfo(sln.GlobalSections, c.prn)
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
