package canvas

import (
	"image"
	"image/color"
	"image/draw"

	"gonum.org/v1/gonum/mat"
)

// BottomAxis defines an Axis as the bottom Axis of an Axes
const BottomAxis = 0

// LeftAxis defines an Axis as the left Axis of an Axes
const LeftAxis = 1

// TopAxis defines an Axis as the top Axis of an Axes
const TopAxis = 2

// RightAxis defines an Axis as the right Axis of an Axes
const RightAxis = 3

// Axis represents a Primitive for horizontal and vertical axes with Axes as its parent.
type Axis struct {
	Primitive
	W        int
	Min, Max float64
	Type     byte
	Parent   *Axes
	Typer    *Typer
}

// NewAxis creates a new Axis linked to an Axes.
// The parameter location can be set to BottomAxis, LeftAxis, TopAxis or RightAxis.
func NewAxis(parent *Axes, location byte) (*Axis, error) {
	var ax Axis
	var o, s [2]float64

	switch location {
	case BottomAxis:
		o = [2]float64{0, -0.1}
		s = [2]float64{1, 0.2}
	case LeftAxis:
		o = [2]float64{-0.1, 0}
		s = [2]float64{0.2, 1}
	case TopAxis:
		o = [2]float64{0, 0.9}
		s = [2]float64{1, 0.2}
	case RightAxis:
		o = [2]float64{0.9, 0}
		s = [2]float64{0.2, 1}
	}

	ax.W = 2

	ax.Parent = parent
	ax.Type = location
	ax.Origin = o
	ax.Size = s
	Tc := mat.NewDense(3, 3, []float64{
		s[0], 0, o[0],
		0, s[1], o[1],
		0, 0, 1,
	})
	ax.T = append(ax.T, parent.T...)
	ax.T = append(ax.T, Tc)
	ax.FillColor = color.Transparent

	parent.Children = append(parent.Children, &ax)
	return &ax, nil
}

// Render creates a Typer to be used by the children Labels.
// The size of Typer is calculated whenever Axis is requested to render.
// This ensures the size is updated on any parent's change.
func (a *Axis) Render(dst draw.Image) {
	if len(a.Children) == 0 {
		return
	}
	l := a.Children[0].(*Label)
	bounds := l.Bounds()
	height := bounds.Max.Y - bounds.Min.Y
	t, _ := NewTyper(height * 72 / 300)
	t.XAlign = l.XAlign
	t.YAlign = l.YAlign
	a.Typer = t
}

// Labels adds X labels to the Axis with regular spacing.
func (a *Axis) Labels(X []string, padding float64) {
	var spacing = (1 - padding*2) / (float64(len(X)) - 1)

	switch a.Type {
	case BottomAxis:
		for i := range X {
			l, _ := NewLabel(a, padding+spacing*float64(i), a.Size[0]*(0.4), 0.5, X[i])
			l.YAlign = TopAlign
		}
	case LeftAxis:
		for i := range X {
			l, _ := NewLabel(a, a.Size[1]*(0.4), padding+spacing*float64(i), 0.1, X[i])
			l.XAlign = RightAlign
		}
	case TopAxis:
		for i := range X {
			l, _ := NewLabel(a, padding+spacing*float64(i), a.Size[0]*(0.6), 0.5, X[i])
			l.YAlign = BottomAlign
		}
	case RightAxis:
		for i := range X {
			l, _ := NewLabel(a, a.Size[1]*(0.6), padding+spacing*float64(i), 0.1, X[i])
			l.XAlign = LeftAlign
		}
	}
}

func (a *Axis) bounds() image.Rectangle {
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

	//return image.Rect(min(x0, x1), min(y0, y1), max(x0, x1), max(y0, y1))
	return a.Primitive.Bounds()
}
