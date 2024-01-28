package vdom

import (
	"fmt"
	"math/rand"
	"syscall/js"
)

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

func TK(text string, key string, vals ...any) *HTML {
	return T(text, vals...).WithKey(key)
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
