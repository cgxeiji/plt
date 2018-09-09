// Package plt provides functions to plot data as draw.Image struct.
package plt

import (
	"github.com/cgxeiji/plt/canvas"
	"image"
	"image/color"
	"image/draw"
	"log"
)

var (
	white color.Color = color.RGBA{255, 255, 255, 255}
	black color.Color = color.RGBA{0, 0, 0, 255}
	blue  color.Color = color.RGBA{0, 0, 255, 255}

	test int
)

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

func max(s []float64) float64 {
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

// func font() *freetype.Context {
// 	fontBytes, _ := ioutil.ReadFile("luxisr.ttf")
// 	f, _ := freetype.ParseFont(fontBytes)
//
// 	c := freetype.NewContext()
// 	c.SetDPI(300)
// 	c.SetFont(f)
// 	c.SetFontSize(12)
// 	c.SetSrc(image.Black)
//
// 	return c
// }

// Bar creates a draw.Image struct given X and Y slices of []float64.
// X and Y must have the same length.
func Bar(X []string, Y []float64) (draw.Image, error) {
	var w, h int = 1920, 1080

	fig, err := canvas.NewFigure(float64(w), float64(h))
	if err != nil {
		return nil, err
	}

	bg := image.NewRGBA(image.Rect(0, 0, w, h))

	fig.Render(bg)

	ax, _ := canvas.NewAxes(fig, 0.1, 0.1, 0.8, 0.8)
	ax.Render(bg)

	ax.BarPlot(X, Y)

	for _, c := range ax.Children {
		c.Render(bg)
	}

	return bg, nil
}

func show(dst draw.Image, c canvas.Container) {
	draw.Draw(dst, c.Bounds(), &image.Uniform{c.Color()}, image.ZP, draw.Src)
}

func Render(dst draw.Image, c canvas.Container) {
	draw.Draw(dst, c.Bounds(), &image.Uniform{c.Color()}, image.ZP, draw.Src)
}

func showBorder(dst draw.Image, c canvas.Container) {
	border(dst, c.Bounds(), -5, &image.Uniform{black}, image.ZP, draw.Src)
}
