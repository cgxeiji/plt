package canvas

import (
	"image"
	"image/draw"
	"io/ioutil"
	"log"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"gonum.org/v1/gonum/mat"
)

type Label struct {
	Primitive
	Parent *Axis
	Text   string
}

func (l *Label) Render(dst draw.Image) {
	location := l.Bounds().Min
	l.Parent.Typer.Render(dst, location.X, location.Y, l.Text)
}

func NewLabel(parent *Axis, x, y, h float64, text string) (*Label, error) {
	var l Label
	l.Parent = parent
	l.Origin = [2]float64{x, y}
	l.Size = [2]float64{0, h}
	l.XAlign = CenterAlign
	l.YAlign = CenterAlign
	Tc := mat.DenseCopyOf(iM)
	l.T = append(l.T, parent.T...)
	l.T = append(l.T, Tc)
	l.FillColor = colornames.Black
	l.Text = text

	parent.children = append(parent.children, &l)
	return &l, nil
}

var DefaultTyper = NewDefaultTyper()

type Typer struct {
	Drawer         *font.Drawer
	Height         fixed.Int26_6
	XAlign, YAlign Alignment
}

func (t *Typer) Render(dst draw.Image, X, Y int, text string) {
	d := t.Drawer
	d.Dst = dst
	var dX fixed.Int26_6

	// dX needs to be calculated on render because the length of text changes
	switch t.XAlign {
	case CenterAlign:
		dX = fixed.I(X) - d.MeasureString(text)/2
	case RightAlign:
		dX = fixed.I(X) - d.MeasureString(text)
	case LeftAlign:
		dX = fixed.I(X)
	}
	// dY only needs to account for t.Height because Y is calculated by Primitive.Vector()
	dY := fixed.I(Y) + t.Height

	d.Dot = fixed.Point26_6{
		X: dX,
		Y: dY,
	}
	d.DrawString(text)
}

func NewTyper(size int) (*Typer, error) {
	fg := image.Black

	d := &font.Drawer{
		Src: fg,
		Face: truetype.NewFace(defaultFont, &truetype.Options{
			Size:    float64(size),
			DPI:     300,
			Hinting: font.HintingNone,
		}),
	}

	t := &Typer{
		Drawer: d,
		Height: fixed.I(size * 300 / 72),
	}

	return t, nil
}

var defaultFont = parseFont("luxisr.ttf")

func parseFont(file string) *truetype.Font {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Panic(err)
	}

	ttf, err := truetype.Parse(bytes)
	if err != nil {
		log.Panic(err)
	}

	return ttf
}

func NewDefaultTyper() *Typer {
	t, err := NewTyper(8)
	if err != nil {
		log.Fatal(err)
	}
	return t
}
