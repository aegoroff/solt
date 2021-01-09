package validate

import c9s "github.com/aegoroff/godatastruct/collections"

type actioner interface {
	action(path string, items map[string]c9s.StringHashSet)
}
