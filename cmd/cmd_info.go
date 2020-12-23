package cmd

import (
	"fmt"
	c9s "github.com/aegoroff/godatastruct/collections"
	"github.com/akutz/sortfold"
	"github.com/cheynewallace/tabby"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"solt/msvc"
	"solt/solution"
	"strings"
)

type infoCommand struct {
	baseCommand
}

func newInfo(c *conf) *cobra.Command {
	cc := cobraCreator{
		createCmd: func() executor {
			ic := infoCommand{
				newBaseCmd(c),
			}
			return &ic
		},
		c: c,
	}

	cmd := cc.newCobraCommand("in", "info", "Get information about found solutions")
	return cmd
}

func (c *infoCommand) execute() error {
	foldersTree := msvc.ReadSolutionDir(c.sourcesPath, c.fs)

	solutions := msvc.SelectSolutions(foldersTree)

	for _, sol := range solutions {
		sln := sol.Solution

		c.prn.setColor(color.FgGray)
		c.prn.cprint(" %s\n", sol.Path)
		c.prn.resetColor()

		t := tabby.NewCustom(c.prn.twriter())

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

func showProjectsInfo(projects []*solution.Project, p printer) {
	var byType = make(map[string]int)

	for _, p := range projects {
		byType[p.Type]++
	}

	const format = "  %v\t%v\n"

	p.tprint(format, "Project type", "Count")
	p.tprint(format, "------------", "-----")

	for k, v := range byType {
		p.tprint(format, k, v)
	}
	p.flush()
	fmt.Println()
}

func showSectionsInfo(sections []*solution.Section, p printer) {
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

	p.tprint(format, "Configuration")
	p.tprint(format, "------------")

	sortedConfigurations := configurations.Items()
	sortfold.Strings(sortedConfigurations)

	for _, k := range sortedConfigurations {
		p.tprint(format, k)
	}
	p.flush()
	fmt.Println()

	p.tprint(format, "Platform")
	p.tprint(format, "--------")

	sortedPlatforms := platforms.Items()
	sortfold.Strings(sortedPlatforms)

	for _, k := range sortedPlatforms {
		p.tprint(format, k)
	}
	p.flush()
	p.cprint("\n")
}
