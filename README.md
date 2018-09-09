# Data Plotter in Go
I am designing a data plotter with only Go dependencies to learn and study about Go.
The current state of this package in not ready for usage. You may use it at your own discretion.

## Testing the code
To test the code, run:
`go run main/main.go`

or build:
`go build main/main.go`

Then, open your favorite browser and go to:
`localhost:8000`

## Usage
Import this packages as:
```go
import "github.com/cgxeiji/plt"
```

Extra features under:
```go
import "github.com/cgxeiji/plt/canvas"
```

### Basic Bar Chart

```go
x := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"}
y := []float64{1, 1, 2, 4, 1, 10}
plot, err := plt.Bar(x, y)
if err != nil {
	log.Panic(err)
}

png.Encode(w, plot)
```

### Multiple Charts

```go
plot := image.NewRGBA(image.Rect(0, 0, 1080, 1920))

x := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"}
y := []float64{1, 1, 2, 4, 1, 10}

fig, err := canvas.NewFigure(1080, 1920)
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

fig.RenderAll(plot)

png.Encode(w, plot)
```
