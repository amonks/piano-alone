//go:build wasm || js

package vdom

import (
	"sync"
	"syscall/js"
)

type VDOM struct {
	mu   sync.Mutex
	root js.Value
	prev Element
}

type Obj map[string]any

type Element interface {
	Mount(parent js.Value, index int) js.Value
	Update(parent, self js.Value, prev Element) js.Value
	Unmount(parent, self js.Value)

	Attrs() Obj
	Key() string
}

func New(root js.Value) *VDOM {
	return &VDOM{root: root}
}

func (vdom *VDOM) Render(el Element) {
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
