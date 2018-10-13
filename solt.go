package main

import (
	"fmt"
	"github.com/voxelbrain/goptions"
	"log"
	"os"
	"runtime/pprof"
)

const csharpProjectExt = ".csproj"
const cppProjectExt = ".vcxproj"
const csharpCodeFileExt = ".cs"
const solutionFileExt = ".sln"
const packagesConfigFile = "packages.config"

type options struct {
	Path       string `goptions:"-p, --path, obligatory, description='Path to the sources folder'"`
	CpuProfile string `goptions:"-c, --cpuprofile, description='CPU profile file'"`
	Version    bool   `goptions:"--version, description='Print version'"`

	goptions.Verbs

	// Finds files that not included into any project
	LostFiles struct {
		Filter string `goptions:"-f, --filter, description='Files filter. By default .cs files'"`
	} `goptions:"lostfiles"`

	// Finds projects that not included into any solution within sources folder
	LostProjects struct {
	} `goptions:"lostprojects"`

	// Shows nuget packages used within any folder that contains packages.confing file
	Nuget struct {
		Mismatch bool `goptions:"-m, --mismatch, description='Find packages to consolidate i.e. packages with different versions in the same solution'"`
	} `goptions:"nuget"`

	// Shows solutions information
	Info struct {
	} `goptions:"info"`
}

type command func(options) error

var commands = map[goptions.Verbs]command{
	"lostfiles":    lostfilescmd,
	"lostprojects": lostprojectscmd,
	"nuget":        nugetcmd,
	"info":         infocmd,
}

func main() {
	opt := options{}

	err := goptions.Parse(&opt)

	if opt.Version {
		fmt.Printf("solt v%s\n", Version)
		return
	}

	if len(opt.Verbs) == 0 || err != nil {
		goptions.PrintHelp()
		return
	}

	if opt.CpuProfile != "" {
		f, err := os.Create(opt.CpuProfile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if cmd, found := commands[opt.Verbs]; found {
		err := cmd(opt)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
	}
}
