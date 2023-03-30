// Package visualize generate resource graph
package visualize

// Type interface defines visualization display driver
type Type interface {
	Display() (string, error)
}

// Visualize store visualization configuration
type Visualize struct {
	Type Type
}

// NewVisualizer construct new visualizer
func NewVisualizer(displayType Type) *Visualize {
	return &Visualize{
		Type: displayType,
	}
}

// Render render resource graph
func (v *Visualize) Render() (string, error) {
	display, err := v.Type.Display()
	if err != nil {
		return "", err
	}

	return display, nil
}
