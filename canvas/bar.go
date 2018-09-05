package canvas

import (
	"fmt"
	"golang.org/x/image/colornames"
	"gonum.org/v1/gonum/mat"
)

type Bar struct {
	Primitive
	Parent *Axes
	Loc    [2]float64
	Value  float64
}

func (b *Bar) String() string {
	return fmt.Sprintf("%v\n ...Bar {Value: %v, Location: %v}", b.Primitive.String(), b.Value, b.Loc)
}

func NewBar(parent *Axes, dims ...float64) (*Bar, error) {
	var min, max [2]float64
	fmt.Println("bar,", dims)

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
		return &Bar{}, fmt.Errorf("Dimensions not valid")
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
	bar.BG = colornames.Red

	return &bar, nil
}
