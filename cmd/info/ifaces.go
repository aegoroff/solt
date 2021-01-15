package info

import (
	"solt/solution"
)

type sectioner interface {
	allow(*solution.Section) bool
	run(*solution.Section)
}
