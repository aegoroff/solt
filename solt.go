package main

import (
    "fmt"
    "github.com/voxelbrain/goptions"
    "os"
)

const CSharpProjectExt = ".csproj"
const CppProjectExt = ".vcxproj"
const CSharpCodeFileExt = ".cs"
const SolutionFileExt = ".sln"
const PackagesConfingFile = "packages.config"

type options struct {
    Help      goptions.Help `goptions:"-h, --help, description='Show this help'"`
    Verbosity bool          `goptions:"-v, --verbose, description='Be verbose'"`
    Path      string        `goptions:"-p, --path, obligatory, description='Path to the sources folder'"`
    Version   bool          `goptions:"--version, description='Print version'"`

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

type Command func(options) error

var commands = map[goptions.Verbs]Command{
    "lostfiles":    lostfilescmd,
    "lostprojects": lostprojectscmd,
    "nuget":        nugetcmd,
    "info":         infocmd,
}

type walkEntry struct {
    Size   int64
    Parent string
    Name   string
    IsDir  bool
}

func main() {
    opt := options{}

    err := goptions.Parse(&opt)

    if opt.Version {
        fmt.Printf("solt v%s\n", Version)
        return
    }

    if len(opt.Verbs) == 0 || opt.Help || err != nil {
        goptions.PrintHelp()
        return
    }

    if cmd, found := commands[opt.Verbs]; found {
        err := cmd(opt)
        if err != nil {
            fmt.Fprintln(os.Stderr, "error:", err)
            os.Exit(1)
        }
    }
}
