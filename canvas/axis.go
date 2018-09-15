package canvas

import (
	"image"
	"image/color"
	"image/draw"

	"golang.org/x/image/colornames"
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
			NewTick(a, padding+spacing*float64(i), a.Size[0]*(0.4), 0.2, 2)
		}
	case LeftAxis:
		spacing = (1 - padding) / (float64(len(X)) - 1)
		for i := range X {
			l, _ := NewLabel(a, a.Size[1]*(0.4), spacing*float64(i), 0.1, X[i])
			l.XAlign = RightAlign
			NewTick(a, a.Size[1]*(0.4), spacing*float64(i), 0.2, 2)
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

// Tick represents a tick to be drawn on an Axis
type Tick struct {
	Primitive
	W      int
	Parent *Axis
}

// NewTick creates a new Tick linked to an Axis.
func NewTick(parent *Axis, x, y, l float64, w int) (*Tick, error) {
	var t Tick

	t.Parent = parent
	t.Origin = [2]float64{x, y}
	switch parent.Type {
	case BottomAxis:
		t.Size = [2]float64{0, l}
	case LeftAxis:
		t.Size = [2]float64{l, 0}
	}
	t.W = w
	Tc := mat.DenseCopyOf(I)
	t.T = append(t.T, parent.T...)
	t.T = append(t.T, Tc)
	t.FillColor = colornames.Black

	parent.Children = append(parent.Children, &t)
	return &t, nil
}

// Render makes sure Tick's Bounds gets called.
func (t *Tick) Render(dst draw.Image) {
	draw.Draw(dst, t.Bounds(), &image.Uniform{t.Color()}, image.ZP, draw.Over)
}

// Bounds returns a a specific width in pixels.
func (t *Tick) Bounds() image.Rectangle {
	var x0, y0, x1, y1 int

	v := transform(t)

	x0 = int(v.At(0, 0))
	y0 = int(v.At(1, 0))
	x1 = int(v.At(0, 1))
	y1 = int(v.At(1, 1))

	if x0 == x1 {
		x0 -= t.W / 2
		x1 += t.W / 2
	}
	if y0 == y1 {
		y0 -= t.W / 2
		y1 += t.W / 2
	}

	return image.Rect(min(x0, x1), min(y0, y1), max(x0, x1), max(y0, y1))
}
