package pen

import (
	"image"
	"image/color"
	"image/draw"
)

func absMax(a, b int) int {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}

	if a < b {
		return b
	}
	return a
}

type Shape int

const Square Shape = 0
const Circle Shape = 1

type Nib struct {
	Size  int
	Shape Shape
}

func (n *Nib) Mask() *image.Alpha {
	nib := image.NewAlpha(image.Rect(0, 0, n.Size, n.Size))
	w := n.Size

	switch n.Shape {
	case Circle:
		for ix := 0; ix < w; ix++ {
			for iy := 0; iy < w; iy++ {
				xc := ix - w/2
				yc := iy - w/2
				xc *= xc
				yc *= yc
				wc := w * w / 4
				wc -= (xc + yc)
				if wc >= 0 {
					nib.SetAlpha(ix, iy, color.Alpha{255})
				}
			}
		}
	default:
		nib = nil
	}

	return nib
}

// pen.Line(bg, image.Pt(10, 10), image.Pt(100, 90), 10, blue)
func Line(dst draw.Image, sp image.Point, ep image.Point, w int, c color.Color) {
	w2 := w / 2
	if w2 == 0 {
		w2 = 1
	}
	nib := (&Nib{w, Circle}).Mask()

	draw := func() func(x, y int) {
		if nib == nil {
			return func(x, y int) {
				rect := image.Rect(x, y, x+w, y+w)
				rect.Add(image.Pt(-w/2, -w/2))
				draw.Draw(dst, rect, &image.Uniform{c}, image.ZP, draw.Over)
			}
		}
		return func(x, y int) {
			rect := image.Rect(x, y, x+w, y+w)
			rect.Add(image.Pt(-w/2, -w/2))
			draw.DrawMask(dst, rect, &image.Uniform{c}, image.ZP, nib, image.ZP, draw.Over)
		}
	}()

	slope := ep.Sub(sp)

	m := float64(slope.Y) / float64(slope.X)
	b := sp.Y - int(m*float64(sp.X))

	y := sp.Y

	for x := sp.X; x <= ep.X; x++ {
		y = int(m*float64(x)) + b
		draw(x, y)
	}

}
