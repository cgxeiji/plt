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
axes[1].Render(plot)
axes[2].Render(plot)
```
