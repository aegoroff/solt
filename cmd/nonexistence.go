package cmd

import (
	"github.com/spf13/afero"
	"os"
)

type nonexistence interface {
	includes() []string
	each(include string)
}

type nonexist struct {
	incl []string
}

func (n *nonexist) includes() []string { return n.incl }
func (n *nonexist) each(string)        {}

func find(non nonexistence, fs afero.Fs) []string {
	result := []string{}
	for _, include := range non.includes() {
		non.each(include)

		if _, err := fs.Stat(include); os.IsNotExist(err) {
			result = append(result, include)
		}
	}
	return result
}
