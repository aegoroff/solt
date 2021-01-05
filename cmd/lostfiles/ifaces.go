package lostfiles

type exister interface {
	exist(project string, includes []string)
}
