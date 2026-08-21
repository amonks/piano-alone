# assets

The files the pages inline or fetch. Most of what belongs here is
built rather than written — see `.gitignore` for what each artifact
comes from — so this README is also the anchor that keeps the
directory present in a fresh checkout, which is what lets the
`//go:embed files` in `../assets.go` compile before anything has been
generated.
