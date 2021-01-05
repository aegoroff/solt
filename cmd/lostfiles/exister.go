package lostfiles

import "solt/cmd/api"

type exister interface {
	exist(project string, includes []string)
}

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
