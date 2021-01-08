package cmd

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"solt/cmd/out"
	"testing"
)

func Test_Version(t *testing.T) {
	var tests = []struct {
		name string
		cmd  []string
	}{
		{"version", []string{"version"}},
		{"ver", []string{"ver"}},
		{"ver -d", []string{"ver", "-d"}},
		{"ver -d --cpuprofile", []string{"ver", "-d", "--cpuprofile", "cpu"}},
		{"ver -d --memprofile", []string{"ver", "-d", "--memprofile", "mem"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			ass := assert.New(t)
			memfs := afero.NewMemMapFs()
			env := out.NewMemoryEnvironment()

			// Act
			_ = Execute(memfs, env, test.cmd...)

			// Assert
			ass.Contains(env.String(), Version)
		})
	}
}

func Test_Help(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	memfs := afero.NewMemMapFs()
	env := out.NewMemoryEnvironment()

	// Act
	_ = Execute(memfs, env, "")

	// Assert
	ass.Empty(env.String())
}

func Test_Console(t *testing.T) {
	// Arrange
	memfs := afero.NewMemMapFs()
	env := out.NewConsoleEnvironment()

	// Act
	_ = Execute(memfs, env, "ver")

	// Assert
}

func Test_Write_File(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	memfs := afero.NewMemMapFs()
	env := out.NewMemoryEnvironment()
	fp := "/f"

	// Act
	_ = Execute(memfs, env, "ver", "-o", fp)

	// Assert
	ass.Empty(env.String())
	_, err := memfs.Stat(fp)
	ass.NoError(err)
}

func Test_Write_ReadonlyFile(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	memfs := afero.NewMemMapFs()
	env := out.NewMemoryEnvironment()
	fp := "/f"
	ro := afero.NewReadOnlyFs(memfs)

	// Act
	_ = Execute(ro, env, "ver", "-o", fp)

	// Assert
	ass.Empty(env.String())
	_, err := memfs.Stat(fp)
	ass.Error(err)
}
