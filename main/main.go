package main

import (
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

		x := []float64{1, 2, 3, 4, 5, 6}
		y := []float64{1, 1, 2, 4, 1, 10}
		fig, err := canvas.NewFigure(1080, 1920)
		if err != nil {
			log.Panic(err)
		}

		plt.Render(plot, fig)

		ax1, err := canvas.NewAxes(fig, 0.1, 0.55, 0.8, 0.35)
		if err != nil {
			log.Panic(err)
		}
		plt.Render(plot, ax1)

		ax1.BarPlot(x, y)

		for _, c := range ax1.Children {
			plt.Render(plot, c)
		}

		ax2, err := canvas.NewAxes(fig, 0.1, 0.1, 0.8, 0.35)
		if err != nil {
			log.Panic(err)
		}
		plt.Render(plot, ax2)

		y = []float64{1.6, 2.2, 3.4, 0.2, 0, 0.2, 0.5}

		ax2.BarPlot(nil, y)

		for _, c := range ax2.Children {
			plt.Render(plot, c)
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
