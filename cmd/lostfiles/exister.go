package lostfiles

import "solt/cmd/fw"

type nopExister struct{}

func (*nopExister) exist(string, []string) {}

type realExister struct {
	e *fw.Exister
}

func newExister(validate bool, e *fw.Exister) exister {
	if validate {
		return &realExister{e: e}
	}
	return &nopExister{}
}

func (r *realExister) exist(project string, includes []string) {
	r.e.Validate(project, includes)
}
