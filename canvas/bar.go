package canvas

import (
	"fmt"
	"golang.org/x/image/colornames"
	"gonum.org/v1/gonum/mat"
)

// Bar is a struct that contains the information necessary to render a single bar in a chart.
type Bar struct {
	Primitive
	Parent *Axes
	Loc    [2]float64
	Value  float64
}

func (b *Bar) String() string {
	return fmt.Sprintf("%v\n ...Bar {Value: %v, Location: %v}", b.Primitive.String(), b.Value, b.Loc)
}

// NewBar creates a new *Bar struct belonging to a parent Axes.
// NewBar takes a parent *Axes and a dims [4]float64{startX, startY, width, value(heigth)}
func NewBar(parent *Axes, dims ...float64) (*Bar, error) {
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
		return &Bar{}, fmt.Errorf(
			"Error while creating Bar: Dimensions %v of length %v not valid",
			dims, len(dims))
	}

	var bar Bar
	bar.Parent = parent
	bar.Loc = min
	bar.Value = max[1]

	bar.Parent = parent
	bar.Origin = min
	bar.Size = max
	Tc := mat.NewDense(3, 3, []float64{
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
	})
	bar.T = append(bar.T, parent.T...)
	bar.T = append(bar.T, Tc)
	bar.FillColor = colornames.Red

	parent.children = append(parent.children, &bar)
	return &bar, nil
}
