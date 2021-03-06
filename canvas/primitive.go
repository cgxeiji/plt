package canvas

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"gonum.org/v1/gonum/mat"
)

// Alignment defines how to draw each element of the plot.
type Alignment byte

// Constants used to align a Container's side to
// its origin point.
const (
	// Used for X alignment
	LeftAlign Alignment = 0
	// Used for Y alignment
	BottomAlign Alignment = 0
	// Used for X and Y alignment
	CenterAlign Alignment = 1
	// Used for X alignment
	RightAlign Alignment = 2
	// Used for Y alignment
	TopAlign Alignment = 2
)

// transformer is an interface that makes sure the Primitive returns
// a vector transformed into pixels and its transformation matrix.
type transformer interface {
	Vector() mat.Matrix
	Transform() []*mat.Dense
}

// Auxiliary function to pretty print matrices.
func mp(name string, X mat.Matrix) {
	f := mat.Formatted(X, mat.Prefix(" "), mat.Squeeze())
	fmt.Printf("%v \n % v\n", name, f)
}

// Auxiliary function to pretty print matrices.
func ms(X mat.Matrix) fmt.Formatter {
	return mat.Formatted(X, mat.Prefix(" "), mat.Squeeze())
}

// iM defines an identity matrix.
var iM = mat.NewDense(3, 3, []float64{
	1, 0, 0,
	0, 1, 0,
	0, 0, 1})

func transform(t transformer) *mat.Dense {
	v := t.Vector()
	r, c := v.Dims()
	render := mat.NewDense(r, c, nil)

	trans := mat.DenseCopyOf(iM)
	transforms := t.Transform()
	l := len(transforms) - 1
	for i, m := range transforms {
		if i == l {
			break
		}
		trans.Product(trans, m)
	}

	render.Product(trans, v)

	return render
}

// primitive is the building block of the plotter.
// Most elements used on the plotter are derivatives from primitive.
//
// A primitive holds all the information necessary to draw the element
// into an image.
//
// primitive implements Container.
// Any primitive can contain other primitives
// as a slice of Container in children.
type primitive struct {
	Origin, Size           [2]float64
	T                      []*mat.Dense
	FillColor, StrokeColor color.Color
	XAlign, YAlign         Alignment
	children               []Container
}

// Vector returns a mat.Matrix with two point coordinates that define
// the bounding rectangle of the Primitive.
//
// The coordinates system are relative to the Primitive's parent.
func (p *primitive) Vector() mat.Matrix {
	var v []float64
	switch p.XAlign {
	case CenterAlign:
		v = append(v, p.Origin[0]-p.Size[0]/2, p.Origin[0]+p.Size[0]/2)
	case RightAlign:
		v = append(v, p.Origin[0]-p.Size[0], p.Origin[0])
	case LeftAlign:
		v = append(v, p.Origin[0], p.Origin[0]+p.Size[0])
	}
	switch p.YAlign {
	case CenterAlign:
		v = append(v, p.Origin[1]-p.Size[1]/2, p.Origin[1]+p.Size[1]/2)
	case TopAlign:
		v = append(v, p.Origin[1]-p.Size[1], p.Origin[1])
	case BottomAlign:
		v = append(v, p.Origin[1], p.Origin[1]+p.Size[1])
	}
	v = append(v, 1, 1)
	vec := mat.NewDense(3, 2, v)
	// mp("V =", vec)
	return vec
}

// Transform returns the transformation matrix of the Primitive.
func (p *primitive) Transform() []*mat.Dense {
	return p.T
}

// Render draws the Primitive into a draw.Image interface.
func (p *primitive) Render(dst draw.Image) {
	draw.Draw(dst, p.Bounds(), &image.Uniform{p.Color()}, image.ZP, draw.Over)
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a // Happy path
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a // Happy path
}

// Bounds returns the bounds to be rendered out of a Primitive in pixels.
func (p *primitive) Bounds() image.Rectangle {
	var x0, y0, x1, y1 int

	v := transform(p)

	x0 = int(v.At(0, 0))
	y0 = int(v.At(1, 0))
	x1 = int(v.At(0, 1))
	y1 = int(v.At(1, 1))

	// fmt.Println("x0:", x0, "y0:", y0, "x1:", x1, "y1:", y1)

	return image.Rect(min(x0, x1), min(y0, y1), max(x0, x1), max(y0, y1))
}

// Color returns the fill color of a Primitive.
func (p *primitive) Color() color.Color {
	return p.FillColor
}

func (p *primitive) String() string {
	b := p.Bounds()
	return fmt.Sprintf(
		"Primitive {T: %v, Origin: %v (pixels: %v), Size: %v (pixels: %v)}",
		p.T, p.Origin, b.Min, p.Size, b.Size(),
	)
}

// Children returns a slice of Container from the children of a Primitive.
func (p *primitive) Children() []Container {
	return p.children
}

// Container is an interface that allows access to
// Render and a Primitive's children.
type Container interface {
	Render(draw.Image)
	Children() []Container
}
