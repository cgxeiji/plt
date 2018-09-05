package main

import (
	"flag"
	"fmt"
	"github.com/cgxeiji/plt/canvas"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math/rand"
	"net/http"
	"os/exec"
	"runtime"
	"time"
)

var (
	white color.Color = color.RGBA{255, 255, 255, 255}
	black color.Color = color.RGBA{0, 0, 0, 255}
	blue  color.Color = color.RGBA{0, 0, 255, 255}

	test int
)

func Border(dst draw.Image, r image.Rectangle, w int, src image.Image,
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

func plot(X, Y []float64) (draw.Image, error) {
	var w, h int = 640, 480

	figure, err := canvas.NewFigure(float64(w), float64(h))
	if err != nil {
		return nil, err
	}

	bg := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(bg, bg.Bounds(), &image.Uniform{blue}, image.ZP, draw.Src)

	show(bg, figure)

	if len(X) != len(Y) {
		return nil, fmt.Errorf("Dimensions mismatch (X[%v] != Y[%v])", len(X), len(Y))
	}

	ax := figure.NewAxes()
	show(bg, ax)
	n := float64(len(Y))
	var padding float64 = 0.1
	barW := (2.0 - 4.0*padding) / (3*n - 1)
	spaceW := barW / 2.0

	for i, _ := range X {
		bar, _ := canvas.NewBar(ax, padding+barW/2.0+float64(i)*(barW+spaceW), 0, barW, Y[i])
		bar.XAlign = canvas.CenterAlign
		show(bg, bar)
	}

	return bg, nil
}

func readFlags() {
	flag.IntVar(&test, "test", 123, "Testing flag")
}

func show(dst draw.Image, c canvas.Container) {
	draw.Draw(dst, c.Bounds(), &image.Uniform{c.Color()}, image.ZP, draw.Src)
}

func homeH(w http.ResponseWriter, r *http.Request) {
	// figure, err := canvas.NewFigure(400, 300)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// ax := figure.NewAxes()

	// bar, _ := canvas.NewBar(ax, 0.2, 0, 0.1, 0.5)
	// bar2, _ := canvas.NewBar(ax, 0.3, 0, 0.2, 0.7)
	// bar2.BG = colornames.Black
	// bar.XAlign = canvas.CenterAlign
	// bar2.XAlign = canvas.CenterAlign

	// bg := image.NewRGBA(image.Rect(0, 0, 640, 480))

	// figure.Resize(640, 480)

	// log.Println("[Fig]", figure)
	// log.Println("[Axes]", ax)
	// log.Println("[Bar]", bar)

	// draw.Draw(bg, bg.Bounds(), &image.Uniform{blue}, image.ZP, draw.Src)

	// log.Println("[Fig]", figure.Bounds())
	// show(bg, figure)
	// log.Println("[Axes]", ax.Bounds())
	// show(bg, ax)
	// log.Println("[Bar]", bar.Bounds())
	// show(bg, bar)
	// log.Println(ax.Color())

	// rect := bar2.Bounds()
	// Border(bg, rect, 2, &image.Uniform{blue}, rect.Min, draw.Src)

	x := []float64{0.1, 0.2, 0.3, 0.4, 0.5}
	y := []float64{0.1, 0.1, 0.2, 0.4, 0.1}
	bg, _ := plot(x, y)

	png.Encode(w, bg)

}

func main() {
	rand.Seed(int64(time.Now().Second()))
	readFlags()
	flag.Parse()
	log.Println("test variable", test)
	log.Println(runtime.GOOS)

	http.HandleFunc("/", homeH)

	// url := fmt.Sprintf("http://localhost:%d/", 1234)

	log.Println(http.ListenAndServe(":8000", nil))

}

func open(name string) error {
	var (
		cmd  string
		args []string
	)

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}
	args = append(args, name)

	return exec.Command(cmd, args...).Start()
}
