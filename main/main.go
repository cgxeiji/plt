package main

import (
	"fmt"
	"github.com/cgxeiji/plt"
	"github.com/cgxeiji/plt/canvas"
	"image"
	"image/png"
	"log"
	"net/http"
)

func main() {
	log.Println("Serving on 'localhost:8000'")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		plot := image.NewRGBA(image.Rect(0, 0, 1080, 1920))

		x := []float64{10, 20, 30, 40, 50, 60}
		y := []float64{1, 1, 2, 4, 1, 10}
		fig, err := canvas.NewFigure(1080, 1920)
		if err != nil {
			log.Panic(err)
		}

		fig.Render(plot)

		ax1, err := canvas.NewAxes(fig, 0.1, 0.55, 0.8, 0.35)
		if err != nil {
			log.Panic(err)
		}
		ax1.Render(plot)

		ax1.BarPlot(x, y)

		typer, err := canvas.NewTyper()
		if err != nil {
			log.Fatal(err)
		}

		for i, c := range ax1.Children {
			c.Render(plot)
			label, err := canvas.NewLabel(ax1, c.RelativeOrigin()[0], c.RelativeOrigin()[1]-0.03, fmt.Sprint(x[i]))
			if err != nil {
				log.Fatal(err)
			}
			label.Render(plot, typer)
		}

		ax2, err := canvas.NewAxes(fig, 0.1, 0.1, 0.8, 0.35)
		if err != nil {
			log.Panic(err)
		}
		ax2.Render(plot)

		y = []float64{1.6, 2.2, 3.4, 0.2, 0, 0.2, 0.5}

		ax2.BarPlot(nil, y)

		for i, c := range ax2.Children {
			c.Render(plot)
			label, err := canvas.NewLabel(ax2, c.RelativeOrigin()[0], c.RelativeOrigin()[1]-0.03, fmt.Sprint(i))
			if err != nil {
				log.Fatal(err)
			}
			label.Render(plot, typer)
		}

		png.Encode(w, plot)
	})

	http.HandleFunc("/simple", func(w http.ResponseWriter, r *http.Request) {
		x := []float64{1, 2, 3, 4, 5, 6}
		y := []float64{1, 1, 2, 4, 1, 10}
		plot, err := plt.Bar(x, y)
		if err != nil {
			log.Panic(err)
		}

		png.Encode(w, plot)
	})

	log.Println(http.ListenAndServe(":8000", nil))
}
