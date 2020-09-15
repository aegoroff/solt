package cmd

import (
	"fmt"
	"github.com/aegoroff/godatastruct/collections"
	"github.com/akutz/sortfold"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"solt/msvc"
	"solt/solution"
	"strings"
)

func newInfo() *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "in",
		Aliases: []string{"info"},
		Short:   "Get information about found solutions",
		RunE: func(cmd *cobra.Command, args []string) error {
			foldersTree := msvc.ReadSolutionDir(sourcesPath, appFileSystem)

			solutions := msvc.SelectSolutions(foldersTree)

			for _, sol := range solutions {
				sln := sol.Solution

				appPrinter.setColor(color.FgGray)
				appPrinter.cprint(" %s\n", sol.Path)
				appPrinter.resetColor()

				const format = "  %v\t%v\n"

				appPrinter.tprint(format, "Header", sln.Header)
				appPrinter.tprint(format, "Product", sln.Comment)
				appPrinter.tprint(format, "Visial Studion Version", sln.VisualStudioVersion)
				appPrinter.tprint(format, "Minimum Visial Studion Version", sln.MinimumVisualStudioVersion)

				appPrinter.flush()

				fmt.Println()

				showProjectsInfo(sln.Projects)
				showSectionsInfo(sln.GlobalSections)
			}

			return nil
		},
	}
	return cmd
}

func showProjectsInfo(projects []*solution.Project) {
	var byType = make(map[string]int)

	for _, p := range projects {
		byType[p.Type]++
	}

	const format = "  %v\t%v\n"

	appPrinter.tprint(format, "Project type", "Count")
	appPrinter.tprint(format, "------------", "-----")

	for k, v := range byType {
		appPrinter.tprint(format, k, v)
	}
	appPrinter.flush()
	fmt.Println()
}

func showSectionsInfo(sections []*solution.Section) {
	var configurations = make(collections.StringHashSet)
	var platforms = make(collections.StringHashSet)

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

	appPrinter.tprint(format, "Configuration")
	appPrinter.tprint(format, "------------")

	sortedConfigurations := configurations.Items()
	sortfold.Strings(sortedConfigurations)

	for _, k := range sortedConfigurations {
		appPrinter.tprint(format, k)
	}
	appPrinter.flush()
	fmt.Println()

	appPrinter.tprint(format, "Platform")
	appPrinter.tprint(format, "--------")

	sortedPlatforms := platforms.Items()
	sortfold.Strings(sortedPlatforms)

	for _, k := range sortedPlatforms {
		appPrinter.tprint(format, k)
	}
	appPrinter.flush()
	fmt.Println()
}
