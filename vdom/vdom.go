//go:build wasm || js

package vdom

import (
	"fmt"
	"syscall/js"
)

type VDOM struct {
	root js.Value
	prev *HTML
}

func New(root js.Value) *VDOM {
	return &VDOM{root: root}
}

func (vdom *VDOM) Render(el *HTML) {
	if vdom.prev == nil {
		el.Mount(vdom.root, 0)
	} else {
		self := vdom.root.Get("childNodes").Index(0)
		el.Update(vdom.root, self, vdom.prev)
	}
	vdom.prev = el
}

type Obj map[string]any

type HTML struct {
	Kind     string
	Key      string
	Attrs    Obj
	Children []*HTML
}

func T(text string, vals ...any) *HTML {
	if len(vals) > 0 {
		text = fmt.Sprintf(text, vals...)
	}
	return &HTML{
		Kind:  "TEXT_NODE",
		Attrs: Obj{"value": text},
	}
}

func H(kind string, children ...*HTML) *HTML {
	return &HTML{
		Kind:     kind,
		Children: children,
	}
}

func (html *HTML) WithAttr(k string, v any) *HTML {
	if html.Attrs == nil {
		html.Attrs = Obj{}
	}
	html.Attrs[k] = v
	return html
}

func (html *HTML) WithKey(k string) *HTML {
	html.Key = k
	return html
}

func (html *HTML) HTML() *HTML {
	return html
}

func (html *HTML) Mount(parent js.Value, index int) {
	if html.Kind == "TEXT_NODE" {
		html.MountText(parent, index)
		return
	}

	node := js.Global().Get("document").Call("createElement", html.Kind)
	for k, v := range html.Attrs {
		node.Set(k, v)
	}

	for i, child := range html.Children {
		child.Mount(node, i)
	}

	if index == -1 {
		parent.Call("appendChild", node)
		return
	}

	numSiblings := parent.Get("childNodes").Get("length").Int()
	if index >= numSiblings {
		parent.Call("appendChild", node)
		return
	}

	nextSibling := parent.Get("childNodes").Index(index + 1)
	parent.Call("insertBefore", node, nextSibling)
}

func (html *HTML) Update(parent js.Value, self js.Value, prev *HTML) {
	if html.Kind == "TEXT_NODE" {
		html.UpdateText(parent, self, prev)
		return
	}

	for k := range prev.Attrs {
		if _, stillHas := html.Attrs[k]; !stillHas {
			self.Delete(k)
		}
	}
	for k, v := range html.Attrs {
		prev, had := prev.Attrs[k]
		if !had {
			self.Set(k, v)
		} else if v != prev {
			self.Set(k, v)
		}
	}

	currKeys := map[string]*HTML{}
	for _, h := range html.Children {
		if h.Key != "" {
			currKeys[h.Key] = h
		}
	}

	// Go through the dom and decide what to do with each node: either
	// update it or delete it.
	type update struct {
		prev *HTML
		node js.Value
	}
	toUpdate := map[string]update{}
	toUnmount := map[*HTML]js.Value{}
	for i, h := range prev.Children {
		node := self.Get("childNodes").Index(i)
		if _, retained := currKeys[h.Key]; retained {
			toUpdate[h.Key] = update{prev: h, node: node}
		} else {
			toUnmount[h] = node
		}
	}

	// Delete dom nodes we don't need anymore.
	for _, node := range toUnmount {
		self.Call("removeChild", node)
	}

	// Go through the HTML we'd like to create, build or update it, then
	// append it to the DOM.
	for i, h := range html.Children {
		if up, isUpdate := toUpdate[h.Key]; isUpdate {
			h.Update(self, up.node, up.prev)
		} else {
			h.Mount(self, i)
		}
	}
}

func (html *HTML) Unmount(parent, self js.Value) {
	if html.Kind == "TEXT_NODE" {
		html.UnmountText(parent, self)
		return
	}

	parent.Call("removeChild", self)
}

func (html *HTML) MountText(parent js.Value, index int) {
	parent.Set("innerText", html.Attrs["value"])
}

func (html *HTML) UpdateText(parent js.Value, self js.Value, prev *HTML) {
	parent.Set("innerText", html.Attrs["value"])
}

func (html *HTML) UnmountText(parent, self js.Value) {
	parent.Set("innerText", "")
}
