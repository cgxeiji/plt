package canvas

import (
	"fmt"
	"image"
	"image/draw"
	"log"

	"golang.org/x/image/colornames"
	"gonum.org/v1/gonum/mat"
)

// Axes represents a Primitive with Figure as its parent.
type Axes struct {
	Primitive
	Parent *Figure
}

// NewAxes creates a new Axes linked to a parent Figure.
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

func minSlice(s []float64) float64 {
	if len(s) <= 0 {
		log.Panic("max(s) on an empty slice")
	}
	var m = s[0]
	for _, v := range s {
		if v < m {
			m = v
		}
	}

	return m
}

func maxSlice(s []float64) float64 {
	if len(s) <= 0 {
		log.Panic("max(s) on an empty slice")
	}
	var m = s[0]
	for _, v := range s {
		if v > m {
			m = v
		}
	}

	return m
}

// BarPlot creates a Bar chart inside Axes with X labels and Y values.
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
	var padding = 0.1
	barW := (2.0 - 4.0*padding) / (3*n - 1)
	spaceW := barW / 2.0

	for i := range Y {
		bar, err := NewBar(ax,
			padding+barW/2.0+float64(i)*(barW+spaceW),
			0,
			barW,
			Y[i]/maxY)
		if err != nil {
			return err
		}
		bar.XAlign = CenterAlign
	}

	axX, _ := NewAxis(ax, BottomAxis)

	if X != nil {
		axX.Labels(X, padding+spaceW)
	}

	maxY = maxSlice(Y)
	labelsY := []string{}
	if Y != nil {
		step := (maxY) / float64(4)
		for i := 0; i < 5; i++ {
			labelsY = append(labelsY, fmt.Sprintf("%.2f", step*float64(i)))
		}
	}

	axY, err := NewAxis(ax, LeftAxis)
	if err != nil {
		return err
	}
	axY.Labels(labelsY, 0.1)

	return nil
}

func vmap(value, fmin, fmax, tmin, tmax float64) float64 {
	return (value-fmin)/(fmax-fmin)*(tmax-tmin) + tmin
}

// ScatterPlot creates a Scatter chart inside Axes with X and Y values.
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
	labels := []string{}
	labelsY := []string{}

	var padding = 0.1

	for i := range Y {
		_, err := NewScatterPoint(ax, vmap(X[i], 0, maxX, padding, 1-padding), Y[i]/maxY)
		if err != nil {
			return err
		}
	}

	if X != nil {
		min := minSlice(X)
		step := (maxX - min) / float64(len(X))
		for i := range X {
			labels = append(labels, fmt.Sprintf("%.2f", min+step*float64(i)))
		}
	}

	axX, err := NewAxis(ax, BottomAxis)
	if err != nil {
		return err
	}
	axX.Labels(labels, padding)

	axX2, err := NewAxis(ax, TopAxis)
	if err != nil {
		return err
	}
	axX2.Labels(labels, padding)

	maxY = maxSlice(Y)
	if Y != nil {
		min := minSlice(Y)
		step := (maxY - min) / float64(4)
		for i := 0; i < 5; i++ {
			labelsY = append(labelsY, fmt.Sprintf("%.2f", min+step*float64(i)))
		}
	}

	axY, err := NewAxis(ax, LeftAxis)
	if err != nil {
		return err
	}
	axY.Labels(labelsY, 0.1)

	axY2, err := NewAxis(ax, RightAxis)
	if err != nil {
		return err
	}
	axY2.Labels(labelsY, 0)

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

// Render draws the Axes' border on top of drawing its contents.
func (ax *Axes) Render(dst draw.Image) {
	ax.Primitive.Render(dst)
	border(dst, ax.Bounds(), -2, &image.Uniform{colornames.Black}, image.ZP, draw.Src)
}
