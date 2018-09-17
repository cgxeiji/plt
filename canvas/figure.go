package canvas

import (
	"golang.org/x/image/colornames"
	"gonum.org/v1/gonum/mat"
)

// Figure defines the basic area to draw all the elements of
// the plot.
// This is the top parent container.
type Figure struct {
	primitive
}

// Resize changes the width and height of the Figure.
// It also updates the transformation matrix of Figure.
func (f *Figure) Resize(w, h float64) {
	f.Size = [2]float64{w, h}
	f.T[1].Set(0, 0, w)
	f.T[1].Set(1, 1, h)

	f.T[0].Set(1, 2, h)
}

// NewAxes attaches a new Axes into the Figure.
func (f *Figure) NewAxes() *Axes {
	axes, _ := f.SubAxes(1, 1)
	return axes[0]
}

// SubAxes attaches multiple Axes defined by the number of rows and columns.
func (f *Figure) SubAxes(rows, cols int) ([]*Axes, error) {
	var axes []*Axes

	n := float64(cols)
	var padX, padY float64 = 0.12, 0.08

	spaceW := padX
	axW := (1 - 2*padX - (n-1)*spaceW) / n

	n = float64(rows)
	spaceH := padY
	axH := (1 - 2*padY - (n-1)*spaceH) / n

	for j := rows - 1; j >= 0; j-- {
		for i := 0; i < cols; i++ {
			ax, err := NewAxes(f, padX+float64(i)*(axW+spaceW), padY+float64(j)*(axH+spaceH), axW, axH)
			if err != nil {
				return nil, err
			}
			axes = append(axes, ax)
		}
	}

	return axes, nil
}

// NewFigure creates a new *Figure with width and height in pixels.
func NewFigure(w, h int) (*Figure, error) {
	max := [2]float64{float64(w), float64(h)}

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

	var fig Figure
	fig.T = append(fig.T, T, Tc)
	fig.Resize(max[0], max[1])
	fig.FillColor = colornames.Gainsboro

	return &fig, nil
}
