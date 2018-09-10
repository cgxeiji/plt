// Package plt provides functions to plot data as draw.Image struct.
package plt

import (
	"github.com/cgxeiji/plt/canvas"
	"image"
	"image/color"
	"image/draw"
)

var (
	white color.Color = color.RGBA{255, 255, 255, 255}
	black color.Color = color.RGBA{0, 0, 0, 255}
	blue  color.Color = color.RGBA{0, 0, 255, 255}

	test int
)

// Bar creates a draw.Image struct given X and Y slices of []float64.
// X and Y must have the same length.
func Bar(X []string, Y []float64) (draw.Image, error) {
	var w, h int = 1920, 1080

	fig, err := canvas.NewFigure(float64(w), float64(h))
	if err != nil {
		return nil, err
	}

	ax, _ := canvas.NewAxes(fig, 0.1, 0.1, 0.8, 0.8)
	ax.BarPlot(X, Y)

	plot := Render(fig)
	return plot, nil
}

func Render(f *canvas.Figure) draw.Image {
	dst := image.NewRGBA(f.Bounds())

	renderAll(f, dst)

	return dst
}

func renderAll(c canvas.Container, dst draw.Image) {
	c.Render(dst)
	for _, child := range c.GetChildren() {
		renderAll(child, dst)
	}
}
