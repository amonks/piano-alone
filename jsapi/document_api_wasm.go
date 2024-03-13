package jsapi

import "syscall/js"

func Doc() Document {
	return Document(js.Global().Get("document"))
}
