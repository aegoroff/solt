package msvc

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_readerSolution_readUnexist(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	memfs := afero.NewMemMapFs()
	ch := make(chan *Folder, 4)

	rm := readerSolution{fs: memfs}

	// Act
	rm.read("ddd", ch)
	close(ch)

	// Assert
	f := <-ch
	ass.Nil(f)
}

func Test_readerSolution_readBad(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	memfs := afero.NewMemMapFs()
	dir := "a/"
	_ = afero.WriteFile(memfs, dir+"a.sln", []byte("xxx"), 0644)
	rm := readerSolution{fs: memfs}
	ch := make(chan *Folder, 4)

	// Act
	rm.read(dir, ch)
	close(ch)

	// Assert
	f := <-ch
	ass.Nil(f)
}

func Test_readerMsbuild_readBad(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	memfs := afero.NewMemMapFs()
	path := "a/a.csproj"
	_ = afero.WriteFile(memfs, path, []byte("xxx"), 0644)
	rm := readerMsbuild{fs: memfs}
	ch := make(chan *Folder, 4)

	// Act
	rm.read(path, ch)
	close(ch)

	// Assert
	f := <-ch
	ass.Nil(f)
}

func Test_readerPackagesConfig_readBad(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	memfs := afero.NewMemMapFs()
	path := "a/packages.config"
	_ = afero.WriteFile(memfs, path, []byte("xxx"), 0644)
	rm := readerPackagesConfig{fs: memfs}
	ch := make(chan *Folder, 4)

	// Act
	rm.read(path, ch)
	close(ch)

	// Assert
	f := <-ch
	ass.Nil(f)
}
