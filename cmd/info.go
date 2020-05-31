package cmd

import (
	"fmt"
	"github.com/aegoroff/godatastruct/collections"
	"github.com/aegoroff/godatastruct/rbtree"
	"solt/solution"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:     "in",
	Aliases: []string{"info"},
	Short:   "Get information about found solutions",
	RunE: func(cmd *cobra.Command, args []string) error {
		foldersTree := readProjectDir(sourcesPath, appFileSystem, func(we *walkEntry) {})

		foldersTree.Ascend(func(c rbtree.Node) bool {
			folder := c.Key().(*folder)
			content := folder.content

			for _, solution := range content.solutions {
				sln := solution.solution

				fmt.Printf(" %s\n", solution.path)

				const format = "  %v\t%v\n"
				tw := new(tabwriter.Writer).Init(appWriter, 0, 8, 4, ' ', 0)

				_, _ = fmt.Fprintf(tw, format, "Header", sln.Header)
				_, _ = fmt.Fprintf(tw, format, "Product", sln.Comment)
				_, _ = fmt.Fprintf(tw, format, "Visial Studion Version", sln.VisualStudioVersion)
				_, _ = fmt.Fprintf(tw, format, "Minimum Visial Studion Version", sln.MinimumVisualStudioVersion)

				_ = tw.Flush()

				fmt.Println()

				showProjectsInfo(sln.Projects)
				showSectionsInfo(sln.GlobalSections)
			}

			return true
		})

		if showMemUsage {
			printMemUsage(appWriter)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

func showProjectsInfo(projects []*solution.Project) {
	var byType = make(map[string]int)

	for _, p := range projects {
		byType[p.Type]++
	}

	const format = "  %v\t%v\n"
	tw := new(tabwriter.Writer).Init(appWriter, 0, 8, 4, ' ', 0)

	_, _ = fmt.Fprintf(tw, format, "Project type", "Count")
	_, _ = fmt.Fprintf(tw, format, "------------", "-----")

	for k, v := range byType {
		_, _ = fmt.Fprintf(tw, format, k, v)
	}
	_ = tw.Flush()
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
	tw := new(tabwriter.Writer).Init(appWriter, 0, 8, 4, ' ', 0)

	_, _ = fmt.Fprintf(tw, format, "Configuration")
	_, _ = fmt.Fprintf(tw, format, "------------")

	sortedConfigurations := configurations.Items()
	sort.Strings(sortedConfigurations)

	for _, k := range sortedConfigurations {
		_, _ = fmt.Fprintf(tw, format, k)
	}
	_ = tw.Flush()
	fmt.Println()

	_, _ = fmt.Fprintf(tw, format, "Platform")
	_, _ = fmt.Fprintf(tw, format, "--------")

	sortedPlatforms := platforms.Items()
	sort.Strings(sortedPlatforms)

	for _, k := range sortedPlatforms {
		_, _ = fmt.Fprintf(tw, format, k)
	}
	_ = tw.Flush()
	fmt.Println()
}
