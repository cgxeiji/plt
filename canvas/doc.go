// Package canvas holds all the components required to plot a chart.
//
// All canvas elements except Figure depend on a parent.
// To create a figure, use the method:
//  fig, err := canvas.NewFigure(w, h)
//  if err != nil {
//		panic(err)
//  }
//
// All other elements should be created from a parent
// as shown in this family tree:
//  Figure
//   |- Axes (figure.NewAxes(), figure.SubAxes(c, r))
//       |- Bar Chart (axes.BarPlot(X, Y))
//       |- Scatter Point Chart (axes.ScatterPlot(X, Y))
//
// Canvas uses a primitive as the building block of the plotter.
// A primitive implements Container and holds all the information
// necessary to draw an element into an image.
// Most elements used on the plotter are derivatives from primitive.
// For that reason, they also implement Container.
//
// Any primitive can contain other primitives
// as a slice of Container in children.
package canvas
