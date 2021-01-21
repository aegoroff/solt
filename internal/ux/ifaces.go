package ux

// Displayer provides interface to component that
// displays something
type Displayer interface {
	// Display does display using Tabler specified
	Display(*Tabler)
}
