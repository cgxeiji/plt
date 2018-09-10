package main

import (
	"fmt"
	"github.com/cgxeiji/plt"
	"image/png"
	"log"
	"net/http"
	"time"
)

func main() {
	log.Println("Serving on 'localhost:8000'")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		x := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"}
		y := []float64{1, 1, 2, 4, 1, 10}

		fig, err := plt.Figure(1080, 1920)
		if err != nil {
			log.Panic(err)
		}

		axes, err := fig.SubAxes(3, 1)
		if err != nil {
			log.Panic(err)
		}

		axes[0].BarPlot(x, y)

		sx := []float64{0.0, 1.1, 2.0, 3.0, 3.5, 4.3, 5.0}
		sy := []float64{1.6, 2.2, 3.4, 0.2, 0.0, 0.2, 0.5}

		axes[1].BarPlot(nil, sy)

		axes[2].ScatterPlot(sx, sy)

		plot := plt.Render(fig)

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
