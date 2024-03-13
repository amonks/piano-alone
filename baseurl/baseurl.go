package baseurl

import (
	"fmt"
	"strings"
)

type BaseURL string

func (url BaseURL) Host() string {
	return string(url)
}

func (url BaseURL) Rest(path string) string {
	return url.Host() + pathName(path)
}

func (url BaseURL) WS(path string) string {
	var base string
	switch true {
	case strings.HasPrefix(url.Host(), "https://"):
		base = strings.Replace(url.Host(), "https://", "wss://", 1)
	case strings.HasPrefix(url.Host(), "http://"):
		base = strings.Replace(url.Host(), "http://", "ws://", 1)
	default:
		panic(fmt.Errorf("invalid base url: %s", url))
	}
	return base + pathName(path)
}

func From(url string) BaseURL {
	return BaseURL(url)
}

func pathName(route string) string {
	switch true {
	case strings.HasPrefix(route, "GET "):
		return strings.TrimPrefix(route, "GET ")
	case strings.HasPrefix(route, "POST "):
		return strings.TrimPrefix(route, "POST ")
	default:
		return route
	}
}
