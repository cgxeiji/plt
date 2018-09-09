package canvas

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"gonum.org/v1/gonum/mat"
	"image"
	"image/draw"
	"io/ioutil"
)

type Label struct {
	Primitive
	Parent *Axes
	Text   string
}

func (l *Label) Render(dst draw.Image, typeFont *Typer) {
	bounds := l.Bounds()
	location := bounds.Min
	typeFont.Render(dst, location.X, location.Y, l.Text)
}

func NewLabel(parent *Axes, x, y float64, text string) (*Label, error) {
	var l Label
	l.Parent = parent
	l.Origin = [2]float64{x, y}
	l.XAlign = CenterAlign
	l.YAlign = CenterAlign
	Tc := mat.DenseCopyOf(I)
	l.T = append(l.T, parent.T...)
	l.T = append(l.T, Tc)
	l.FillColor = colornames.Black
	l.Text = text

	return &l, nil
}

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
