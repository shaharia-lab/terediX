package visualize

type Type interface {
	Display() (string, error)
}

type Visualize struct {
	Type Type
}

func NewVisualizer(displayType Type) *Visualize {
	return &Visualize{
		Type: displayType,
	}
}

func (v *Visualize) Render() (string, error) {
	display, err := v.Type.Display()
	if err != nil {
		return "", err
	}

	return display, nil
}
