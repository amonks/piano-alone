package baseurl

import "syscall/js"

func Discover() BaseURL {
	loc := js.Global().Get("document").Get("location")
	protocol := loc.Get("protocol").String()
	host := loc.Get("host").String()
	return BaseURL(protocol + "//" + host)
}
