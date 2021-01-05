package msvc

import (
	"os"
	"path/filepath"
	"strings"
)

type msbuildStandardPaths struct {
	folder string
}

func newMsbuildStandardPaths(folder string) *msbuildStandardPaths {
	b := strings.Builder{}
	b.WriteString(folder)
	b.WriteRune(os.PathSeparator)
	base := b.String()

	return &msbuildStandardPaths{folder: base}
}

func (sp *msbuildStandardPaths) decorate(s string) string {
	r := strings.ReplaceAll(s, "$(MSBuildProjectDirectory)", sp.folder)
	path := strings.ReplaceAll(r, "$(MSBuildThisFileDirectory)", sp.folder)
	if strings.HasPrefix(path, sp.folder) {
		return filepath.Clean(path)
	}
	return filepath.Join(sp.folder, path)
}
