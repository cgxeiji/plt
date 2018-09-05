package canvas

import (
	"fmt"
	"golang.org/x/image/colornames"
	"gonum.org/v1/gonum/mat"
)

type Figure struct {
	Primitive
}

func (f *Figure) Resize(w, h float64) {
	f.Size = [2]float64{w, h}
	f.T[1].Set(0, 0, w)
	f.T[1].Set(1, 1, h)

	f.T[0].Set(1, 2, h)
}

func (f *Figure) NewAxes() *Axes {
	ax, _ := NewAxes(f, 0.1, 0.1, 0.8, 0.8)
	return &ax
}

func NewFigure(dims ...float64) (*Figure, error) {
	var min, max [2]float64

	switch l := len(dims); l {
	case 0:
		min = [2]float64{0, 0}
		max = [2]float64{640, 480}
	case 1:
		min = [2]float64{0, 0}
		max = [2]float64{dims[0], dims[0]}
	case 2:
		min = [2]float64{0, 0}
		max = [2]float64{dims[0], dims[1]}
	default:
		return &Figure{}, fmt.Errorf("Dimensions not valid")
	}

	var fig Figure
	fig.Origin = min
	T := mat.NewDense(3, 3, []float64{
		1, 0, 0,
		0, -1, max[1],
		0, 0, 1,
	})
	Tc := mat.NewDense(3, 3, []float64{
		max[0], 0, 0,
		0, max[1], 0,
		0, 0, 1,
	})
	fig.T = append(fig.T, T, Tc)
	fig.Resize(max[0], max[1])
	fig.BG = colornames.Gray

	return &fig, nil
}
