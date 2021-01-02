package sys

import (
	"os"
	"strings"
)

// ToValidPath creates valid OS specific path from path specified
func ToValidPath(p string) string {
	if os.PathSeparator != '/' {
		return p
	}
	return strings.ReplaceAll(p, "\\", "/")
}
