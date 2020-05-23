solt
====

[![Build status](https://ci.appveyor.com/api/projects/status/tgx6ai9erbgfq2ij?svg=true)](https://ci.appveyor.com/project/aegoroff/solt) [![codecov](https://codecov.io/gh/aegoroff/solt/branch/master/graph/badge.svg)](https://codecov.io/gh/aegoroff/solt) [![Go Report Card](https://goreportcard.com/badge/github.com/aegoroff/solt)](https://goreportcard.com/report/github.com/aegoroff/solt)

**SOL**ution **T**ool is a small commandline app written in Go that allows you to easily analyze
sources and Microsoft Visual Studion solutions and projects.
The tool can find files that aren't included into any project and projects that
are not included into any solution. Additionally the tool shows some useful
solution statistic

Command line syntax:
--------------------
```
SOLution Tool that analyzes Microsoft Visual Studio solutions and projects

Usage:
  solt [flags]
  solt [command]

Available Commands:
  help         Help about any command
  info         Get information about found solutions
  lostfiles    Find lost files in the folder specified
  lostprojects Find projects that not included into any solution
  nuget        Get nuget packages information within projects or find Nuget mismatches in solution
  version      Print the version number of solt

Flags:
  -h, --help          help for solt
  -p, --path string   REQUIRED. Path to the sources folder

Use "solt [command] --help" for more information about a command.
```
### Search lost files syntax:

```
Find lost files in the folder specified

Usage:
  solt lostfiles [flags]

Aliases:
  lostfiles, lf

Flags:
  -f, --file string   Lost files filter extension. If not set .cs extension used (default ".cs")
  -h, --help          help for lostfiles
  -l, --onlylost      Show only lost files. Don't show unexist files. If not set all shown
  -r, --remove        Remove lost files

Global Flags:
  -p, --path string   REQUIRED. Path to the sources folder
```
### Search lost projects syntax:
```
Find projects that not included into any solution

Usage:
  solt lostprojects [flags]

Aliases:
  lostprojects, lp

Flags:
  -h, --help   help for lostprojects

Global Flags:
  -p, --path string   REQUIRED. Path to the sources folder
```
### Nuget information syntax:
```
Get nuget packages information within projects or find Nuget mismatches in solution

Usage:
  solt nuget [flags]

Aliases:
  nuget, nu

Flags:
  -h, --help       help for nuget
  -m, --mismatch   Find packages to consolidate i.e. packages with different versions in the same solution

Global Flags:
  -p, --path string   REQUIRED. Path to the sources folder
```