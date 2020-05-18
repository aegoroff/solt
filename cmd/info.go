package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"solt/solution"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"in"},
	Short:   "Get information about found solutions",
	RunE: func(cmd *cobra.Command, args []string) error {
		var solutions []string
		readProjectDir(sourcesPath, func(we *walkEntry) {
			ext := strings.ToLower(filepath.Ext(we.Name))
			if ext == solutionFileExt {
				solutions = append(solutions, filepath.Join(we.Parent, we.Name))
			}
		})

		for _, sol := range solutions {
			sln, err := solution.Parse(sol)

			if err != nil {
				continue
			}

			fmt.Printf(" %s\n", sol)

			const format = "  %v\t%v\n"
			tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 4, ' ', 0)

			fmt.Fprintf(tw, format, "Header", sln.Header)
			fmt.Fprintf(tw, format, "Product", sln.Comment)
			fmt.Fprintf(tw, format, "Visial Studion Version", sln.VisualStudioVersion)
			fmt.Fprintf(tw, format, "Minimum Visial Studion Version", sln.MinimumVisualStudioVersion)

			tw.Flush()

			fmt.Println()

			showProjectsInfo(sln.Projects)
			showSectionsInfo(sln.GlobalSections)

			if err != nil {
				return err
			}
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
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 4, ' ', 0)

	fmt.Fprintf(tw, format, "Project type", "Count")
	fmt.Fprintf(tw, format, "------------", "-----")

	for k, v := range byType {
		fmt.Fprintf(tw, format, k, v)
	}
	tw.Flush()
	fmt.Println()
}

func showSectionsInfo(sections []*solution.Section) {
	var configurations = make(map[string]bool)
	var platforms = make(map[string]bool)

	for _, s := range sections {
		if s.Name != "SolutionConfigurationPlatforms" {
			continue
		}
		for _, item := range s.Items {
			parts := strings.Split(item.Key, "|")
			configuration := parts[0]
			platform := parts[1]
			if _, ok := configurations[configuration]; !ok {
				configurations[configuration] = true
			}
			if _, ok := platforms[platform]; !ok {
				platforms[platform] = true
			}
		}
	}

	const format = "  %v\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 4, ' ', 0)

	fmt.Fprintf(tw, format, "Configuration")
	fmt.Fprintf(tw, format, "------------")

	for k := range configurations {
		fmt.Fprintf(tw, format, k)
	}
	tw.Flush()
	fmt.Println()

	fmt.Fprintf(tw, format, "Platform")
	fmt.Fprintf(tw, format, "--------")

	for k := range platforms {
		fmt.Fprintf(tw, format, k)
	}
	tw.Flush()
	fmt.Println()
}
