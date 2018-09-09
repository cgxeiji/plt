package main

import (
	"fmt"
	"github.com/cgxeiji/plt"
	"github.com/cgxeiji/plt/canvas"
	"image"
	"image/png"
	"log"
	"net/http"
	"time"
)

func main() {
	log.Println("Serving on 'localhost:8000'")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		plot := image.NewRGBA(image.Rect(0, 0, 1080, 1920))

		x := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"}
		y := []float64{1, 1, 2, 4, 1, 10}
		fig, err := canvas.NewFigure(1080, 1920)
		if err != nil {
			log.Panic(err)
		}

		fig.Render(plot)

		axes, err := fig.SubAxes(3, 1)
		if err != nil {
			log.Panic(err)
		}
		axes[0].Render(plot)

		axes[0].BarPlot(x, y)

		for _, c := range axes[0].Children {
			c.Render(plot)
		}

		axes[1].Render(plot)

		sx := []float64{0.0, 1.1, 2.0, 3.0, 3.5, 4.3, 5.0}
		sy := []float64{1.6, 2.2, 3.4, 0.2, 0.0, 0.2, 0.5}

		axes[1].BarPlot(nil, sy)
		for _, c := range axes[1].Children {
			c.Render(plot)
		}

		axes[2].Render(plot)
		axes[2].ScatterPlot(sx, sy)
		for _, c := range axes[2].Children {
			c.Render(plot)
		}

		png.Encode(w, plot)

		fmt.Println("Rendered in:", time.Now().Sub(startTime))
	})

	http.HandleFunc("/simple", func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		x := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"}
		y := []float64{1, 1, 2, 4, 1, 10}
		plot, err := plt.Bar(x, y)
		if err != nil {
			log.Panic(err)
		}

		png.Encode(w, plot)
		fmt.Println("Rendered in:", time.Now().Sub(startTime))
	})

	log.Println(http.ListenAndServe(":8000", nil))
}
