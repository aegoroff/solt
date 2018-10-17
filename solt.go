package main

import (
	"github.com/urfave/cli"
	"log"
	"os"
)

const csharpProjectExt = ".csproj"
const cppProjectExt = ".vcxproj"
const solutionFileExt = ".sln"
const packagesConfigFile = "packages.config"

var sourcesPath string
var lostFilesFilter string
var findNugetMismatches bool

func main() {
	app := cli.NewApp()
	app.Name = "solt"
	app.Version = Version
	app.Usage = "SOLution Tool that analyzes Microsoft Visual Studio solutions and projects"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "path, p",
			Usage:       "REQUIRED. Path to the sources folder",
			Destination: &sourcesPath,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "lostfiles",
			Aliases: []string{"f"},
			Usage:   "Find lost files in the folder specified",
			Action: func(c *cli.Context) error {
				return actionLoader(c, lostfilescmd)
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "file, f",
					Value:       ".cs",
					Usage:       "Lost files filter extension",
					Destination: &lostFilesFilter,
				},
			},
		},
		{
			Name:    "lostprojects",
			Aliases: []string{"p"},
			Usage:   "Find projects that not included into any solution",
			Action: func(c *cli.Context) error {
				return actionLoader(c, lostprojectscmd)
			},
		},
		{
			Name:    "nuget",
			Aliases: []string{"n"},
			Usage:   "Get nuget packages information within projects or find Nuget mismatches in solution",
			Action: func(c *cli.Context) error {
				return actionLoader(c, nugetcmd)
			},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:        "mismatch, m",
					Usage:       "Find packages to consolidate i.e. packages with different versions in the same solution",
					Destination: &findNugetMismatches,
				},
			},
		},
		{
			Name:    "info",
			Aliases: []string{"i"},
			Usage:   "Get information about found solutions",
			Action: func(c *cli.Context) error {
				return actionLoader(c, infocmd)
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func actionLoader(c *cli.Context, action func(*cli.Context) error) error {
	if len(sourcesPath) == 0 {
		cli.ShowAppHelp(c)
		return nil
	}
	return action(c)
}
