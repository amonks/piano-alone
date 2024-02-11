//go:build tools

package main

import (
	_ "golang.org/x/tools/cmd/stringer"
	_ "golang.org/x/vuln/cmd/govulncheck"
	_ "honnef.co/go/tools/cmd/staticcheck"
)
