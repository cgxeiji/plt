package canvas

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"gonum.org/v1/gonum/mat"
)

const LeftAlign byte = 0
const BottomAlign byte = 0

const CenterAlign byte = 1

const RightAlign byte = 2
const TopAlign byte = 2

type Renderer interface {
	Vector() mat.Matrix
	Transform() []*mat.Dense
}

func mp(name string, X mat.Matrix) {
	f := mat.Formatted(X, mat.Prefix(" "), mat.Squeeze())
	fmt.Printf("%v \n % v\n", name, f)
}

func ms(X mat.Matrix) fmt.Formatter {
	return mat.Formatted(X, mat.Prefix(" "), mat.Squeeze())
}

var I *mat.Dense = mat.NewDense(3, 3, []float64{
	1, 0, 0,
	0, 1, 0,
	0, 0, 1})

func transform(r Renderer) *mat.Dense {
	v := r.Vector()
	rows, cols := v.Dims()
	render := mat.NewDense(rows, cols, nil)

	trans := mat.DenseCopyOf(I)
	transforms := r.Transform()
	max_len := len(transforms) - 1
	for i, t := range transforms {
		if i == max_len {
			break
		}
		// h := fmt.Sprintf("T%v =", i)
		// mp(h, t)
		trans.Product(trans, t)
	}

	render.Product(trans, v)
	// mp("T =", trans)
	// mp("R = T * V", render)

	return render
}

type Primitive struct {
	Origin, Size           [2]float64
	T                      []*mat.Dense
	FillColor, StrokeColor color.Color
	XAlign, YAlign         byte
	Children               []Container
}

func (p *Primitive) Vector() mat.Matrix {
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

func (p *Primitive) Transform() []*mat.Dense {
	return p.T
}

func (p *Primitive) Render(dst draw.Image) {
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

func (p *Primitive) Bounds() image.Rectangle {
	var x0, y0, x1, y1 int

	v := transform(p)

	x0 = int(v.At(0, 0))
	y0 = int(v.At(1, 0))
	x1 = int(v.At(0, 1))
	y1 = int(v.At(1, 1))

	// fmt.Println("x0:", x0, "y0:", y0, "x1:", x1, "y1:", y1)

	return image.Rect(min(x0, x1), min(y0, y1), max(x0, x1), max(y0, y1))
}

func (p *Primitive) Color() color.Color {
	return p.FillColor
}

func (p *Primitive) String() string {
	b := p.Bounds()
	return fmt.Sprintf("Primitive {T: %v, Origin: %v (pixels: %v), Size: %v (pixels: %v)}", p.T, p.Origin, b.Min, p.Size, b.Size())
}

func (p *Primitive) GetChildren() []Container {
	return p.Children
}

type Container interface {
	Bounds() image.Rectangle
	Color() color.Color
	Render(draw.Image)
	GetChildren() []Container
}
