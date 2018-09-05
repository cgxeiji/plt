package main

import (
	"github.com/cgxeiji/plt"
	"image/png"
	"log"
	"net/http"
)

func main() {
	log.Println("Serving on 'localhost:8000'")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
