solt
====

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/b8b9bdf73cfb4e97888b6ff7b48bfc84)](https://app.codacy.com/manual/egoroff/solt?utm_source=github.com&utm_medium=referral&utm_content=aegoroff/solt&utm_campaign=Badge_Grade_Dashboard)
[![Build status](https://ci.appveyor.com/api/projects/status/tgx6ai9erbgfq2ij?svg=true)](https://ci.appveyor.com/project/aegoroff/solt) [![codecov](https://codecov.io/gh/aegoroff/solt/branch/master/graph/badge.svg)](https://codecov.io/gh/aegoroff/solt) [![Go Report Card](https://goreportcard.com/badge/github.com/aegoroff/solt)](https://goreportcard.com/report/github.com/aegoroff/solt)

**SOL**ution **T**ool is a small commandline app written in Go that allows you to easily analyze
sources and Microsoft Visual Studion solutions and projects.
The tool can find files that aren't included into any project and projects that
are not included into any solution. Additionally the tool shows some useful
solution statistic

## Install the pre-compiled binary

**homebrew** (only on macOS and Linux for now):

Add my tap (do it once):
```sh
brew tap aegoroff/tap
```
And the install solt:
```sh
brew install solt
```
Update solt if already installed:
```sh
brew upgrade solt
```

**manually**:

Download the pre-compiled binaries from the [releases](https://github.com/aegoroff/solt/releases) and
copy to the desired location.

Command line syntax:
--------------------
```
SOLution Tool that analyzes Microsoft Visual Studio solutions and projects

Usage:
  solt [flags]
  solt [command]

Available Commands:
  help        Help about any command
  in          Get information about found solutions
  lf          Find lost files in the folder specified
  lp          Find projects that not included into any solution
  nu          Get nuget packages information within solutions,
              projects or find Nuget mismatches in solution
  va          Validates SDK projects within solution(s)
  ver         Print the version number of solt

Flags:
      --cpuprofile string   Runs CPU profiling if --diag option set.
                            If not set profiling not started.
                            Correct file path should be set here
  -d, --diag                Show application diagnostic after run
  -h, --help                help for solt
      --memprofile string   Runs memory profiling if --diag option set.
                            If not set profiling not started.
                            Correct file path should be set here
  -p, --path string         REQUIRED. Path to the sources folder

Use "solt [command] --help" for more information about a command.
```
### Search lost files syntax:

```
Find lost files in the folder specified

Usage:
  solt lf [flags]

Aliases:
  lf, lostfiles

Flags:
  -a, --all           Search all lost files including that have 
                      links to but not exists in file system
  -f, --file string   Lost files filter extension. (default ".cs")
  -h, --help          help for lf
  -r, --remove        Remove lost files

Global Flags:
      --cpuprofile string   Runs CPU profiling if --diag option set.
                            If not set profiling not started.
                            Correct file path should be set here
  -d, --diag                Show application diagnostic after run
      --memprofile string   Runs memory profiling if --diag option set.
                            If not set profiling not started.
                            Correct file path should be set here
  -p, --path string         REQUIRED. Path to the sources folder
```
### Search lost projects syntax:
```
Find projects that not included into any solution

Usage:
  solt lp [flags]

Aliases:
  lp, lostprojects

Flags:
  -h, --help   help for lp

Global Flags:
      --cpuprofile string   Runs CPU profiling if --diag option set.
                            If not set profiling not started.
                            Correct file path should be set here
  -d, --diag                Show application diagnostic after run
      --memprofile string   Runs memory profiling if --diag option set.
                            If not set profiling not started.
                            Correct file path should be set here
  -p, --path string         REQUIRED. Path to the sources folder

```
### Nuget information syntax:
```
Get nuget packages information within solutions

Usage:
  solt nu [flags]
  solt nu [command]

Aliases:
  nu, nuget

Available Commands:
  p           Get nuget packages information by projects' folders 
              i.e. from packages.config or SDK project files

Flags:
  -h, --help       help for nu
  -m, --mismatch   Find packages to consolidate i.e. packages 
                   with different versions in the same solution

Global Flags:
      --cpuprofile string   Runs CPU profiling if --diag option set.
                            If not set profiling not started. 
                            Correct file path should be set here
  -d, --diag                Show application diagnostic after run
      --memprofile string   Runs memory profiling if --diag option set. 
                            If not set profiling not started. 
                            Correct file path should be set here
  -p, --path string         REQUIRED. Path to the sources folder
```
### Validate SDK projects syntax:
```
Validates SDK projects within solution(s)

Usage:
  solt va [flags]
  solt va [command]

Aliases:
  va, validate

Available Commands:
  fix         Fixes redundant SDK projects references

Flags:
  -h, --help   help for va

Global Flags:
      --cpuprofile string   Runs CPU profiling if --diag option set.
                            If not set profiling not started.
                            Correct file path should be set here
  -d, --diag                Show application diagnostic after run
      --memprofile string   Runs memory profiling if --diag option set.
                            If not set profiling not started.
                            Correct file path should be set here
  -p, --path string         REQUIRED. Path to the sources folder
```