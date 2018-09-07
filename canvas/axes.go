package canvas

import (
	"fmt"
	"golang.org/x/image/colornames"
	"gonum.org/v1/gonum/mat"
	"image"
	"log"
)

type Axes struct {
	Primitive
	Parent   *Figure
	Children []Container
}

func NewAxes(parent *Figure, dims ...float64) (*Axes, error) {
	var o, s [2]float64

	switch l := len(dims); l {
	case 4:
		o = [2]float64{dims[0], dims[1]}
		s = [2]float64{dims[2], dims[3]}
	default:
		return &Axes{}, fmt.Errorf("Dimensions not valid")
	}

	var ax Axes
	ax.Parent = parent
	ax.Origin = o
	ax.Size = s
	Tc := mat.NewDense(3, 3, []float64{
		s[0], 0, o[0],
		0, s[1], o[1],
		0, 0, 1,
	})
	ax.T = append(ax.T, parent.T...)
	ax.T = append(ax.T, Tc)
	ax.FillColor = colornames.White

	return &ax, nil
}

func maxSlice(s []float64) float64 {
	if len(s) <= 0 {
		log.Panic("max(s) on an empty slice")
	}
	var m float64 = s[0]
	for _, v := range s {
		if v > m {
			m = v
		}
	}

	return m
}

func (ax *Axes) BarPlot(X, Y []float64) error {
	if X != nil {
		if len(X) != len(Y) {
			return fmt.Errorf(
				"Dimensions mismatch (X[%v] != Y[%v])",
				len(X), len(Y))
		}
	}

	maxY := maxSlice(Y) / 0.9

	n := float64(len(Y))
	var padding float64 = 0.1
	barW := (2.0 - 4.0*padding) / (3*n - 1)
	spaceW := barW / 2.0

	for i, _ := range Y {
		bar, err := NewBar(ax,
			padding+barW/2.0+float64(i)*(barW+spaceW),
			0,
			barW,
			Y[i]/maxY)
		if err != nil {
			return err
		}
		bar.XAlign = CenterAlign
		ax.Children = append(ax.Children, bar)
	}

	return nil
}

type Axis struct {
	Primitive
	W      int
	Parent *Axes
}

func NewAxis(parent *Axes, which byte) (*Axis, error) {
	var ax Axis

	ax.Origin = [2]float64{0, 0}
	switch which {
	case 0:
		ax.Size = [2]float64{0, 1}
		ax.XAlign = RightAlign
	case 1:
		ax.Size = [2]float64{1, 0}
		ax.YAlign = TopAlign
	}

	ax.W = 8

	ax.Parent = parent
	Tc := I
	ax.T = append(ax.T, parent.T...)
	ax.T = append(ax.T, Tc)
	ax.FillColor = colornames.Black

	return &ax, nil
}

func (a *Axis) Bounds() image.Rectangle {
	var x0, y0, x1, y1 int

	v := transform(a)

	x0 = int(v.At(0, 0))
	y0 = int(v.At(1, 0))
	x1 = int(v.At(0, 1))
	y1 = int(v.At(1, 1))

	if x0 == x1 {
		x0 -= a.W
		y0 += a.W * 2
	}
	if y0 == y1 {
		y1 += a.W
		x0 -= a.W * 2
	}

	return image.Rect(min(x0, x1), min(y0, y1), max(x0, x1), max(y0, y1))

}
