package canvas

import (
	"fmt"

	"golang.org/x/image/colornames"
	"gonum.org/v1/gonum/mat"
)

// bar is a struct that contains the information necessary to render a single bar in a chart.
type bar struct {
	primitive
	Parent *Axes
	Loc    [2]float64
	Value  float64
}

func (b *bar) String() string {
	return fmt.Sprintf("%v\n ...Bar {Value: %v, Location: %v}", b.primitive.String(), b.Value, b.Loc)
}

// newBar creates a new *Bar struct belonging to a parent Axes.
// newBar takes a parent *Axes and a dims [4]float64{startX, startY, width, value(heigth)}
func newBar(parent *Axes, dims ...float64) (*bar, error) {
	var min, max [2]float64

	switch l := len(dims); l {
	case 0:
		min = [2]float64{0, 0}
		max = [2]float64{0.1, 0.1}
	case 1:
		min = [2]float64{0, 0}
		max = [2]float64{dims[0], dims[0]}
	case 2:
		min = [2]float64{0, 0}
		max = [2]float64{dims[0], dims[1]}
	case 4:
		min = [2]float64{dims[0], dims[1]}
		max = [2]float64{dims[2], dims[3]}
	default:
		return &bar{}, fmt.Errorf(
			"Error while creating Bar: Dimensions %v of length %v not valid",
			dims, len(dims))
	}

	var b bar
	b.Parent = parent
	b.Loc = min
	b.Value = max[1]

	b.Parent = parent
	b.Origin = min
	b.Size = max
	Tc := mat.NewDense(3, 3, []float64{
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
	})
	b.T = append(b.T, parent.T...)
	b.T = append(b.T, Tc)
	b.FillColor = colornames.Red

	parent.children = append(parent.children, &b)
	return &b, nil
}
