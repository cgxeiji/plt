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

		x := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"}
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

		for _, c := range ax1.Children {
			c.Render(plot)
		}

		ax2, err := canvas.NewAxes(fig, 0.1, 0.1, 0.8, 0.35)
		if err != nil {
			log.Panic(err)
		}
		ax2.Render(plot)

		y = []float64{1.6, 2.2, 3.4, 0.2, 0, 0.2, 0.5}

		ax2.BarPlot(nil, y)

		for _, c := range ax2.Children {
			c.Render(plot)
		}

		png.Encode(w, plot)
	})

	http.HandleFunc("/simple", func(w http.ResponseWriter, r *http.Request) {
		x := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"}
		y := []float64{1, 1, 2, 4, 1, 10}
		plot, err := plt.Bar(x, y)
		if err != nil {
			log.Panic(err)
		}

		png.Encode(w, plot)
	})

	log.Println(http.ListenAndServe(":8000", nil))
}
