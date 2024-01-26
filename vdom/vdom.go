//go:build wasm || js

package vdom

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"syscall/js"
)

type VDOM struct {
	mu   sync.Mutex
	root js.Value
	prev *HTML
}

func New(root js.Value) *VDOM {
	return &VDOM{root: root}
}

func (vdom *VDOM) Render(el *HTML) {
	vdom.mu.Lock()
	defer vdom.mu.Unlock()
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
		Key:   "TEXT_NODE",
		Kind:  "TEXT_NODE",
		Attrs: Obj{"value": text},
	}
}

func TK(text string, key string, vals ...any) *HTML {
	return T(text, vals...).WithKey(key)
}

func H(kind string, children ...*HTML) *HTML {
	return &HTML{
		Kind:     kind,
		Key:      kind,
		Children: children,
	}
}

func HK(kind string, key string, children ...*HTML) *HTML {
	return H(kind, children...).WithKey(key)
}

func (html *HTML) WithAttr(k string, v any) *HTML {
	if html.Attrs == nil {
		html.Attrs = Obj{}
	}
	html.Attrs[k] = v
	return html
}

func (html *HTML) WithKey(k string) *HTML {
	html.Key = html.Kind + "." + k
	return html
}

func (html *HTML) HTML() *HTML {
	return html
}

func (html *HTML) Mount(parent js.Value, index int) js.Value {
	if html.Kind == "TEXT_NODE" {
		return html.MountText(parent, index)
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
	vdom        *HTML
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

func (html *HTML) Update(parent js.Value, self js.Value, prev *HTML) js.Value {

	if html.Kind == "TEXT_NODE" {
		return html.UpdateText(parent, self, prev)
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
		case string:
			// only update string attrs if they changed
			if v != prev.Attrs[k] {
				self.Set(k, v)
			}
		default:
			// if we don't know whether an attr is comparable, just
			// go ahead and try to update it.
			self.Set(k, v)
		}
	}

	domChildren := make(mountedChildren, len(prev.Children))
	domChildrenByKey := map[string]*mountedChild{}
	childNodes := self.Get("childNodes")
	for i, c := range prev.Children {
		mc := &mountedChild{
			parent:     self,
			vdom:       c,
			node:       childNodes.Index(i),
			isRetained: false,
		}
		if !mc.node.Truthy() {
			panic("falsy node")
		}
		domChildrenByKey[c.Key] = mc
		domChildren[i] = mc
	}

	keyset := map[string]struct{}{}
	for i, c := range html.Children {
		if _, dupe := keyset[c.Key]; dupe {
			panic(fmt.Errorf("reused key in '%s': '%s'", html.Key, c.Key))
		}
		keyset[c.Key] = struct{}{}

		if mc, hasMC := domChildrenByKey[c.Key]; hasMC {
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
		domChildrenByKey[c.Key] = mc
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
	if html.Kind == "TEXT_NODE" {
		html.UnmountText(parent, self)
		return
	}

	parent.Call("removeChild", self)
}

func (html *HTML) MountText(parent js.Value, index int) js.Value {
	node := js.Global().Get("document").Call("createTextNode", html.Attrs["value"])
	sibs := parent.Get("childNodes")
	if index == -1 || index >= sibs.Length()-1 {
		parent.Call("appendChild", node)
	} else {
		parent.Call("insertBefore", node, sibs.Index(index+1))
	}
	return node
}

func (html *HTML) UpdateText(parent js.Value, self js.Value, prev *HTML) js.Value {
	new, old := html.Attrs["value"], prev.Attrs["value"]
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
