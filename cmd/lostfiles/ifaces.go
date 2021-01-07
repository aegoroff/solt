package lostfiles

import "solt/cmd/fw"

type exister interface {
	exist(project string, includes []string)
	print(p fw.Printer)
}
