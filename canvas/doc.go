// Package canvas holds all the components required to plot a chart.
// It can be imported as:
//  import github.com/cgxeiji/plt/canvas
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
