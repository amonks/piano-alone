// go:bulid js && wasm

package canvas

import (
	"syscall/js"

	"monks.co/piano-alone/c2d"
)

type SceneNode interface {
	Draw(c2d c2d.C2D, bounds Bounds)
}

var _ SceneNode = SceneNodeFunc(nil)

type SceneNodeFunc func(c2d c2d.C2D, bounds Bounds)

func (f SceneNodeFunc) Draw(c2d c2d.C2D, bounds Bounds) {
	f(c2d, bounds)
}

type Bounds struct {
	X, Y, Width, Height float64
}

func (b Bounds) Proportional(p Bounds) Bounds {
	return Bounds{
		X:      b.X + p.X*b.Width,
		Y:      b.Y + p.Y*b.Height,
		Width:  b.Width * p.Width,
		Height: b.Height * p.Height,
	}
}

func (b Bounds) Intersection(other Bounds) Bounds {
	return Bounds{
		X:      b.X + other.X,
		Y:      b.Y + other.Y,
		Width:  min(b.Width, other.Width),
		Height: min(b.Height, other.Height),
	}
}

func Draw(canvas js.Value, nodes []SceneNode) js.Value {
	bounds := canvas.Call("getBoundingClientRect")
	w, h := bounds.Get("width").Float(), bounds.Get("height").Float()
	canvas.Set("width", w)
	canvas.Set("height", h)

	c2d := c2d.C2D(canvas.Call("getContext", "2d"))
	c2d.ClearRect(0, 0, w, h)
	for _, n := range nodes {
		n.Draw(c2d, Bounds{0, 0, w, h})
	}
	return canvas
}
