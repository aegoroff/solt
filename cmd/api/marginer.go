package api

import "strings"

// Marginer makes string with margins
type Marginer struct {
	value int
}

// NewMarginer creates new Marginer instance using margin value
func NewMarginer(value int) *Marginer {
	return &Marginer{value: value}
}

// Margin creates new string with margin
func (m *Marginer) Margin(s string) string {
	sb := strings.Builder{}
	for i := 0; i < m.value; i++ {
		sb.WriteString(" ")
	}
	sb.WriteString(s)

	return sb.String()
}
