package vdom

import (
	"syscall/js"

	"monks.co/piano-alone/c2d"
)

type Canvas struct {
	key   string
	attrs Obj
	nodes []SceneNode
}

var _ Element = &Canvas{}

type SceneNode interface {
	Draw(c2d.C2D)
}

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
	return node
}

func (c *Canvas) Update(parent, self js.Value, prev Element) js.Value {
	c2d := c2d.C2D(self.Call("getContext", "2d"))
	c2d.ClearRect(0, 0, self.Get("width").Float(), self.Get("height").Float())
	for _, n := range c.nodes {
		n.Draw(c2d)
	}
	return self
}

func (c *Canvas) Unmount(parent, self js.Value) {
	parent.Call("removeChild", self)
}
