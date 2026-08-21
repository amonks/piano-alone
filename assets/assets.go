// Package assets holds the files the pages inline or fetch: the
// compiled stylesheet, the wasm player and the Go runtime's loader
// shim, and the script that starts them.
//
// The embed is of the directory rather than of each file by name,
// because most of what it holds is generated: a named embed of a file
// the build has not produced yet fails to compile, which would mean a
// fresh checkout could not build until it had run the generators. A
// missing file surfaces at render time instead, as an error on the one
// page that wanted it.
package assets

import (
	"embed"
	"io/fs"
)

//go:embed files
var embedded embed.FS

// FS is the asset directory.
var FS = mustSub(embedded, "files")

func mustSub(f fs.FS, dir string) fs.FS {
	sub, err := fs.Sub(f, dir)
	if err != nil {
		panic(err)
	}
	return sub
}

// Read returns one asset's bytes.
func Read(name string) ([]byte, error) { return fs.ReadFile(FS, name) }
