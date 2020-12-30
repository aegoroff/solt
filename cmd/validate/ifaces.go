package validate

import "github.com/aegoroff/godatastruct/collections"

type actioner interface {
	action(path string, items map[string]collections.StringHashSet)
}
