// Package plt provides functions to plot data as draw.Image struct.
package plt

import (
	"image"
	"image/draw"

	"github.com/cgxeiji/plt/canvas"
)

// Bar creates a draw.Image struct given X and Y slices
// of []string and []float64 respectively.
// X and Y must have the same length.
func Bar(X []string, Y []float64) (draw.Image, error) {
	var w, h int = 1920, 1080

	fig, err := canvas.NewFigure(w, h)
	if err != nil {
		return nil, err
	}

	ax, _ := canvas.NewAxes(fig, 0.1, 0.1, 0.8, 0.8)
	ax.BarPlot(X, Y)

	plot := Render(fig)
	return plot, nil
}

// Figure returns a *canvas.Figure struct with a size of
// w and h in pixels (int).
func Figure(w, h int) (*canvas.Figure, error) {
	return canvas.NewFigure(w, h)
}

// Render draws a Figure with all its children into a draw.Image interface.
func Render(f *canvas.Figure) draw.Image {
	dst := image.NewRGBA(f.Bounds())

	renderAll(f, dst)

	return dst
}

func renderAll(c canvas.Container, dst draw.Image) {
	c.Render(dst)
	for _, child := range c.Children() {
		renderAll(child, dst)
	}
}
