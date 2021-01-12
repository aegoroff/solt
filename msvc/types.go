package msvc

import (
	"github.com/aegoroff/dirstat/scan"
	"github.com/spf13/afero"
)

const (
	// SolutionFileExt defines visual studio extension
	SolutionFileExt    = ".sln"
	csharpProjectExt   = ".csproj"
	cppProjectExt      = ".vcxproj"
	packagesConfigFile = "packages.config"
)

// ProjectHandler defines project handler function prototype
type ProjectHandler func(*MsbuildProject, *Folder)

// StringDecorator defines string decorating function
type StringDecorator func(s string) string

type filesystem struct{ fs afero.Fs }

func newFs(fs afero.Fs) scan.Filesystem {
	return &filesystem{fs: fs}
}

func (f *filesystem) Open(path string) (scan.File, error) {
	return f.fs.Open(path)
}
