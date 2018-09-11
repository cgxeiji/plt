package canvas

import (
	"fmt"
	"image"
	"image/draw"
	"log"

	"golang.org/x/image/colornames"
	"gonum.org/v1/gonum/mat"
)

type Axes struct {
	Primitive
	Parent *Figure
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

	parent.Children = append(parent.Children, &ax)

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

func (ax *Axes) BarPlot(X []string, Y []float64) error {
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

		if X != nil {
			_, err := NewLabel(ax, bar.Origin[0], -0.1, 0.08, X[i])
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	return nil
}

func vmap(value, fmin, fmax, tmin, tmax float64) float64 {
	return (value-fmin)/(fmax-fmin)*(tmax-tmin) + tmin
}

func (ax *Axes) ScatterPlot(X, Y []float64) error {
	if X != nil {
		if len(X) != len(Y) {
			return fmt.Errorf(
				"Dimensions mismatch (X[%v] != Y[%v])",
				len(X), len(Y))
		}
	}

	maxY := maxSlice(Y) / 0.9
	maxX := maxSlice(X)

	var padding float64 = 0.1

	for i, _ := range Y {
		point, err := NewScatterPoint(ax, vmap(X[i], 0, maxX, padding, 1-padding), Y[i]/maxY)
		if err != nil {
			return err
		}

		if X != nil {
			_, err := NewLabel(ax, point.Origin[0], -0.1, 0.08, fmt.Sprint(X[i]))
			if err != nil {
				log.Fatal(err)
			}

		}
	}

	return nil
}

func border(dst draw.Image, r image.Rectangle, w int, src image.Image,
	sp image.Point, op draw.Op) {
	// inside r
	if w > 0 {
		// top
		draw.Draw(dst, image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Min.Y+w), src, sp, op)
		// left
		draw.Draw(dst, image.Rect(r.Min.X, r.Min.Y+w, r.Min.X+w, r.Max.Y-w),
			src, sp.Add(image.Pt(0, w)), op)
		// right
		draw.Draw(dst, image.Rect(r.Max.X-w, r.Min.Y+w, r.Max.X, r.Max.Y-w),
			src, sp.Add(image.Pt(r.Dx()-w, w)), op)
		// bottom
		draw.Draw(dst, image.Rect(r.Min.X, r.Max.Y-w, r.Max.X, r.Max.Y),
			src, sp.Add(image.Pt(0, r.Dy()-w)), op)
		return
	}

	// outside r;
	w = -w
	// top
	draw.Draw(dst, image.Rect(r.Min.X-w, r.Min.Y-w, r.Max.X+w, r.Min.Y),
		src, sp.Add(image.Pt(-w, -w)), op)
	// left
	draw.Draw(dst, image.Rect(r.Min.X-w, r.Min.Y, r.Min.X, r.Max.Y), src,
		sp.Add(image.Pt(-w, 0)), op)
	// right
	draw.Draw(dst, image.Rect(r.Max.X, r.Min.Y, r.Max.X+w, r.Max.Y), src,
		sp.Add(image.Pt(r.Dx(), 0)), op)
	// bottom
	draw.Draw(dst, image.Rect(r.Min.X-w, r.Max.Y, r.Max.X+w, r.Max.Y+w),
		src, sp.Add(image.Pt(-w, 0)), op)
}

func (ax *Axes) Render(dst draw.Image) {
	ax.Primitive.Render(dst)
	border(dst, ax.Bounds(), -2, &image.Uniform{colornames.Black}, image.ZP, draw.Src)
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
