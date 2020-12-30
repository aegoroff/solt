package api

import "strings"

// Marginer makes string with margins
type Marginer struct {
	value int
	char  string
}

// NewMarginer creates new Marginer instance using margin value and space as margin char
func NewMarginer(value int) *Marginer {
	return NewCustomMarginer(value, " ")
}

// NewCustomMarginer creates new Marginer instance using margin value and specified as margin char
func NewCustomMarginer(value int, c string) *Marginer {
	return &Marginer{value: value, char: c}
}

// Margin creates new string with margin
func (m *Marginer) Margin(s string) string {
	sb := strings.Builder{}
	for i := 0; i < m.value; i++ {
		sb.WriteString(m.char)
	}
	sb.WriteString(s)

	return sb.String()
}

// NewUnderline creates new dashes string that can be used as underline
// the length of the line is equal len of string specified
func NewUnderline(s string) string {
	return NewCustomMarginer(len(s), "-").Margin("")
}
