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
	Parent *Axes
	Text   string
}

func (l *Label) Render(dst draw.Image) {
	bounds := l.Bounds()
	height := bounds.Max.Y - bounds.Min.Y
	t, _ := NewTyper(height * 72 / 300) // TODO: Change for a faster typer initialization
	location := bounds.Min
	t.Render(dst, location.X, location.Y, l.Text)
}

func NewLabel(parent *Axes, x, y, h float64, text string) (*Label, error) {
	var l Label
	l.Parent = parent
	l.Origin = [2]float64{x, y}
	l.Size = [2]float64{0, h}
	l.XAlign = CenterAlign
	l.YAlign = CenterAlign
	Tc := mat.DenseCopyOf(I)
	l.T = append(l.T, parent.T...)
	l.T = append(l.T, Tc)
	l.FillColor = colornames.Black
	l.Text = text

	parent.Children = append(parent.Children, &l)
	return &l, nil
}

var DefaultTyper = NewDefaultTyper()

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
