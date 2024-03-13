//go:build tools

package main

import (
	_ "github.com/a-h/templ/cmd/templ"
	_ "golang.org/x/tools/cmd/stringer"
	_ "golang.org/x/vuln/cmd/govulncheck"
	_ "honnef.co/go/tools/cmd/staticcheck"
)
