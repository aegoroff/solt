package fw

import "solt/msvc"

// Includer provides includes extractor from msvc.container
type Includer struct {
	includes []string
	exister  Exister
}

// NewIncluder creates new Includer instance
func NewIncluder(exister Exister) *Includer {
	return &Includer{exister: exister}
}

// Solution method that called on each solution while iterating solutions
func (i *Includer) Solution(s *msvc.VisualStudioSolution) {
	i.From(s)
}

// From gets includes from container and validates their existence
func (i *Includer) From(c container) {
	includes := c.Items()
	i.includes = append(i.includes, includes...)

	i.exister.Validate(c.Path(), includes)
}

// Includes gets includes extracted from container
func (i *Includer) Includes() []string {
	return i.includes
}
