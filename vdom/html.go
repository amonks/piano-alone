package vdom

import (
	"fmt"
	"math/rand"
	"sort"
	"syscall/js"
)

type HTML struct {
	kind     string
	key      string
	attrs    Obj
	children []Element
}

func T(text string, vals ...any) *HTML {
	if len(vals) > 0 {
		text = fmt.Sprintf(text, vals...)
	}
	return &HTML{
		key:   "TEXT_NODE",
		kind:  "TEXT_NODE",
		attrs: Obj{"value": text},
	}
}
func H(kind string, children ...Element) Element {
	return &HTML{
		kind:     kind,
		key:      kind,
		children: children,
	}
}

func TK(text string, key string, vals ...any) Element {
	return T(text, vals...).WithKey(key)
}
func HK(kind string, key string, children ...Element) Element {
	return H(kind, children...).WithKey(key)
}

var _ Element = &HTML{}

func (html *HTML) WithAttr(k string, v any) Element {
	if html.attrs == nil {
		html.attrs = Obj{}
	}
	html.attrs[k] = v
	return html
}

func (html *HTML) Attrs() Obj {
	return html.attrs
}

func (html *HTML) WithKey(k string) Element {
	html.key = html.kind + "." + k
	return html
}

func (html *HTML) Key() string {
	return html.key
}

func (html *HTML) Mount(parent js.Value, index int) js.Value {
	if html.kind == "TEXT_NODE" {
		return html.MountText(parent, index)
	}

	node := js.Global().Get("document").Call("createElement", html.kind)
	for k, v := range html.attrs {
		node.Set(k, v)
	}

	for i, child := range html.children {
		child.Mount(node, i)
	}

	if index == -1 {
		parent.Call("appendChild", node)
		return node
	}

	siblingNodes := parent.Get("childNodes")
	numSiblings := siblingNodes.Get("length").Int()
	if index >= numSiblings {
		parent.Call("appendChild", node)
		return node
	}

	nextSibling := siblingNodes.Index(index + 1)
	parent.Call("insertBefore", node, nextSibling)
	return node
}

type mountedChild struct {
	parent      js.Value
	targetIndex int
	vdom        Element
	node        js.Value
	isRetained  bool
}

type mountedChildren []*mountedChild

var _ sort.Interface = mountedChildren{}

func (mc mountedChildren) Len() int           { return len(mc) }
func (mc mountedChildren) Less(i, j int) bool { return mc[i].targetIndex < mc[j].targetIndex }
func (mc mountedChildren) Swap(left, right int) {
	if right < left {
		left, right = right, left
	}
	parent := mc[left].parent
	sibs := parent.Get("childNodes")
	var afterLeft js.Value
	for i := 0; i < sibs.Length(); i++ {
		if sibs.Index(i).Equal(mc[left].node) {
			afterLeft = sibs.Index(i + 1)
			break
		}
	}
	if afterLeft.Equal(mc[right].node) {
		parent.Call("insertBefore", mc[right].node, mc[left].node)
		return
	}
	parent.Call("replaceChild", mc[left].node, mc[right].node)
	parent.Call("insertBefore", mc[right].node, afterLeft)
}

func (html *HTML) Update(parent js.Value, self js.Value, _prev Element) js.Value {
	prev := _prev.(*HTML)

	if html.kind == "TEXT_NODE" {
		return html.UpdateText(parent, self, prev)
	}

	for k := range prev.attrs {
		if _, stillHas := html.attrs[k]; !stillHas {
			self.Delete(k)
		}
	}
	for k, v := range html.attrs {
		switch v.(type) {
		case js.Func:
			// special case: never update func attrs, as they are
			// incomparable
		case string:
			// only update string attrs if they changed
			if v != prev.attrs[k] {
				self.Set(k, v)
			}
		default:
			// if we don't know whether an attr is comparable, just
			// go ahead and try to update it.
			self.Set(k, v)
		}
	}

	domChildren := make(mountedChildren, len(prev.children))
	domChildrenByKey := map[string]*mountedChild{}
	childNodes := self.Get("childNodes")
	for i, c := range prev.children {
		mc := &mountedChild{
			parent:     self,
			vdom:       c,
			node:       childNodes.Index(i),
			isRetained: false,
		}
		if !mc.node.Truthy() {
			panic("falsy node")
		}
		domChildrenByKey[c.Key()] = mc
		domChildren[i] = mc
	}

	keyset := map[string]struct{}{}
	for i, c := range html.children {
		if _, dupe := keyset[c.Key()]; dupe {
			panic(fmt.Errorf("reused key in '%s': '%s'", html.Key, c.Key))
		}
		keyset[c.Key()] = struct{}{}

		if mc, hasMC := domChildrenByKey[c.Key()]; hasMC {
			mc.targetIndex = i
			mc.node = c.Update(self, mc.node, mc.vdom)
			mc.isRetained = true
			continue
		}

		mc := &mountedChild{
			parent:      self,
			node:        c.Mount(self, -1),
			targetIndex: i,
			isRetained:  true,
		}
		domChildren = append(domChildren, mc)
		domChildrenByKey[c.Key()] = mc
	}

	for i := 0; i < len(domChildren); {
		c := domChildren[i]
		if !c.isRetained {
			self.Call("removeChild", c.node)
			domChildren = append(domChildren[:i], domChildren[i+1:]...)
		} else {
			i++
		}
	}

	sort.Sort(domChildren)

	return self
}

func (html *HTML) Unmount(parent, self js.Value) {
	if html.kind == "TEXT_NODE" {
		html.UnmountText(parent, self)
		return
	}

	parent.Call("removeChild", self)
}

func (html *HTML) MountText(parent js.Value, index int) js.Value {
	node := js.Global().Get("document").Call("createTextNode", html.attrs["value"])
	sibs := parent.Get("childNodes")
	if index == -1 || index >= sibs.Length()-1 {
		parent.Call("appendChild", node)
	} else {
		parent.Call("insertBefore", node, sibs.Index(index+1))
	}
	return node
}

func (html *HTML) UpdateText(parent js.Value, self js.Value, prev *HTML) js.Value {
	new, old := html.attrs["value"], prev.attrs["value"]
	if new != old {
		self.Set("data", new)
	}
	return self
}

func (html *HTML) UnmountText(parent, self js.Value) {
	parent.Call("removeChild", self)
}

func randomColor() string {
	hue := rand.Intn(360)
	return fmt.Sprintf("hsl(%ddeg, 50%%, 90%%)", hue)
}