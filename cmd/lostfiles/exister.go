package lostfiles

import "solt/cmd/api"

type nopExister struct{}

func (*nopExister) exist(string, []string) {}

type realExister struct {
	e *api.Exister
}

func newExister(validate bool, e *api.Exister) exister {
	if validate {
		return &realExister{e: e}
	}
	return &nopExister{}
}

func (r *realExister) exist(project string, includes []string) {
	r.e.Validate(project, includes)
}
