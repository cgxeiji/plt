package canvas

import (
	"golang.org/x/image/colornames"
)

type ScatterPoint struct {
	Primitive
	Parent *Axes
	X, Y   float64
}

func NewScatterPoint(parent *Axes, x, y float64) (*ScatterPoint, error) {
	var point ScatterPoint
	point.Parent = parent
	point.X = x
	point.Y = y
	point.Origin = [2]float64{x, y}
	point.Size = [2]float64{0.01, 0.01}
	point.T = append(point.T, parent.T...)
	point.T = append(point.T, nil)

	point.FillColor = colornames.Green

	return &point, nil
}
