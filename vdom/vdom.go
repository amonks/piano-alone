//go:build wasm || js

package vdom

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"syscall/js"
)

const debug = false

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

func H(kind string, key string, children ...*HTML) *HTML {
	return &HTML{
		Kind:     kind,
		Key:      key,
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
	if debug {
		var styles []string
		if s, got := html.Attrs["style"]; got {
			styles = append(styles, s.(string))
		}
		styles = append(styles, fmt.Sprintf("background-color: %s", randomColor()))
		style := strings.Join(styles, "; ")
		node.Set("style", style)
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
	//
	// handle updating self
	//

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
		switch v.(type) {
		case js.Func:
			// special case: never update func attrs, as they are
			// incomparable
			if debug {
				log.Printf("not updating func attr %s", k)
			}
		case string:
			// only update string attrs if they changed
			if v != prev.Attrs[k] {
				if debug {
					log.Printf("updating attr %s", k)
				}
				self.Set(k, v)
			}
		default:
			// if we don't know whether an attr is comparable, just
			// go ahead and try to update it.
			if debug {
				log.Printf("updating attr %s", k)
			}
			self.Set(k, v)
		}
	}

	//
	// handle recurring on children
	//
	// TODO: do a proper "minimal set of edits" sequence comparison algorithm.
	// The standard one is from this paper:
	//     https://publications.mpi-cbg.de/Wu_1990_6334.pdf
	// Here are some implemenattions:
	//     go: https://github.com/cubicdaiya/gonp
	//     js; used for vdom: https://github.com/thi-ng/umbrella/blob/develop/packages/diff/src/array.ts#L53
	//

	// short-circuit: if we have the same sequence of children, skip a
	// bunch of work and just update them all.
	shouldShortCircuit := true
	if len(html.Children) != len(prev.Children) {
		shouldShortCircuit = false
	} else {
		for i, h := range html.Children {
			if prev.Children[i].Key != h.Key {
				shouldShortCircuit = false
				break
			}
		}
	}
	if shouldShortCircuit {
		nodes := self.Get("childNodes")
		for i, h := range html.Children {
			h.Update(self, nodes.Index(i), prev.Children[i])
		}
		return
	}

	// We know that elements need to be created, destroyed, or re-ordered.
	// Let's do some work.

	// Collect the set of current keys.
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

	// Go through the elements that should exist, build or update them as
	// appropriate.
	//
	// BUG: this doesn't always order the children correctly.
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

func randomColor() string {
	hue := rand.Intn(360)
	return fmt.Sprintf("hsl(%ddeg, 50%%, 90%%)", hue)

}
