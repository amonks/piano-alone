package baseurl

import (
	"strings"
)

type BaseURL string

var NoHost = BaseURL("")

func (url BaseURL) Host() string {
	return string(url)
}

func (url BaseURL) Rest(path string, substitutions ...string) string {
	if len(substitutions)%2 != 0 {
		panic("odd number of substitutions")
	}
	p := pathName(path)
	for i := 0; i < len(substitutions); i += 2 {
		k, v := substitutions[i], substitutions[i+1]
		p = strings.Replace(p, "{"+k+"}", v, 1)
	}
	return url.Host() + p
}

func (url BaseURL) WS(path string) string {
	var base string
	switch true {
	case strings.HasPrefix(url.Host(), "https://"):
		base = strings.Replace(url.Host(), "https://", "wss://", 1)
	case strings.HasPrefix(url.Host(), "http://"):
		base = strings.Replace(url.Host(), "http://", "ws://", 1)
	default:
		panic("can't make ws url from hostless baseURL")
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
