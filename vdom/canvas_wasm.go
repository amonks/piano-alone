// go:bulid js && wasm

package vdom

import (
	"syscall/js"

	"monks.co/piano-alone/c2d"
)

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

type Canvas struct {
	c2d           c2d.C2D
	key           string
	attrs         Obj
	nodes         []SceneNode
	width, height float64
}

var _ Element = &Canvas{}

func C(nodes ...SceneNode) *Canvas {
	return &Canvas{nodes: nodes}
}

func (c *Canvas) Attrs() Obj  { return c.attrs }
func (c *Canvas) Key() string { return c.key }

func (c *Canvas) Mount(parent js.Value, index int) js.Value {
	node := js.Global().Get("document").Call("createElement", "canvas")
	for k, v := range c.attrs {
		node.Set(k, v)
	}
	sibs := parent.Get("childNodes")
	if index == -1 || index >= sibs.Length()-1 {
		parent.Call("appendChild", node)
	} else {
		parent.Call("insertBefore", node, sibs.Index(index+1))
	}
	bounds := node.Call("getBoundingClientRect")
	w, h := bounds.Get("width").Float(), bounds.Get("height").Float()
	node.Set("width", w*2)
	node.Set("height", h*2)
	return node
}

func (c *Canvas) Update(parent, self js.Value, prev Element) js.Value {
	c2d := c2d.C2D(self.Call("getContext", "2d"))
	w, h := self.Get("width").Float(), self.Get("height").Float()
	c2d.ClearRect(0, 0, w, h)
	for _, n := range c.nodes {
		n.Draw(c2d, Bounds{0, 0, w, h})
	}
	return self
}

func (c *Canvas) Unmount(parent, self js.Value) {
	parent.Call("removeChild", self)
}

type SceneNode interface {
	Draw(c2d c2d.C2D, bounds Bounds)
}

var _ SceneNode = SceneNodeFunc(nil)

type SceneNodeFunc func(c2d c2d.C2D, bounds Bounds)

func (f SceneNodeFunc) Draw(c2d c2d.C2D, bounds Bounds) {
	f(c2d, bounds)
}
