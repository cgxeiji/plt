package canvas

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/draw"
	"io/ioutil"
)

type Typer struct {
	Drawer *font.Drawer
	Height fixed.Int26_6
}

func (t *Typer) Render(dst draw.Image, X, Y int, text string) {
	d := t.Drawer
	d.Dst = dst
	d.Dot = fixed.Point26_6{
		X: fixed.I(X) - d.MeasureString(text)/2,
		Y: fixed.I(Y) + t.Height/2,
	}
	d.DrawString(text)
}

func NewTyper() (*Typer, error) {
	bytes, err := ioutil.ReadFile("luxisr.ttf")
	if err != nil {
		return nil, err
	}

	ttf, err := truetype.Parse(bytes)
	if err != nil {
		return nil, err
	}

	fg := image.Black

	h := font.HintingNone

	d := &font.Drawer{
		Src: fg,
		Face: truetype.NewFace(ttf, &truetype.Options{
			Size:    12,
			DPI:     300,
			Hinting: h,
		}),
	}

	t := &Typer{
		Drawer: d,
		Height: fixed.I(12 * 300 / 72),
	}

	return t, nil
}
