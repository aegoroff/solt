package main

import (
    "fmt"
    "github.com/voxelbrain/goptions"
    "os"
)

const CSharpProjectExt = ".csproj"
const CSharpCodeFileExt = ".cs"
const SolutionFileExt = ".sln"
const PackagesConfingFile = "packages.config"

type options struct {
    Help      goptions.Help `goptions:"-h, --help, description='Show this help'"`
    Verbosity bool          `goptions:"-v, --verbose, description='Be verbose'"`
    Path      string        `goptions:"-p, --path, obligatory, description='Path to the project'"`

    goptions.Verbs

    Lost struct {
        Filter string `goptions:"-f, --filter, description='Files filter. By default .cs files'"`
    } `goptions:"lost"`

    Info struct {
        Exclude string `goptions:"-e, --exclude, description='Do not include this version into output'"`
    } `goptions:"info"`

    Analyze struct {
        Solution string `goptions:"-s, --solution, description='solution file to analyze'"`
    } `goptions:"analyze"`
}

type Command func(options) error

var commands = map[goptions.Verbs]Command{
    "lost":   lostcmd,
    "info":   infocmd,
    "analyze":   analyzecmd,
}

type walkEntry struct {
    Size   int64
    Parent string
    Name   string
    IsDir  bool
}

func main() {
    opt := options{}

    goptions.ParseAndFail(&opt)

    if len(opt.Verbs) == 0 {
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
