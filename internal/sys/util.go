package sys

import (
	"github.com/spf13/afero"
	"os"
)

// CheckExistence validates files passed to be present in file system
// The list of non exist files returned
func CheckExistence(files []string, fs afero.Fs) []string {
	result := []string{}
	for _, f := range files {
		if _, err := fs.Stat(f); os.IsNotExist(err) {
			result = append(result, f)
		}
	}
	return result
}
