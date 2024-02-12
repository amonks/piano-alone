//go:build js && wasm

package storage

import "syscall/js"

type Storage string

const (
	Local   Storage = "localStorage"
	Session Storage = "sessionStorage"
)

func (s Storage) Get(k string) string {
	got := js.Global().Get(string(s)).Call("getItem", k)
	if !got.Truthy() {
		return ""
	}
	return got.String()
}

func (s Storage) Set(k, v string) {
	js.Global().Get(string(s)).Call("setItem", k, v)
}
