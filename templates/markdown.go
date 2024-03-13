package templates

import (
	"bytes"
	"context"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/a-h/templ"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

func markdown(filename string) templ.Component {
	input, err := os.ReadFile(filepath.Join("templates", filename))
	if err != nil {
		log.Fatalf("could not find markdown file: templates/%s", filename)
	}

	md := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
		goldmark.WithExtensions(
			extension.Typographer,
			extension.Linkify,
		),
	)

	var buf bytes.Buffer
	if err := md.Convert(input, &buf); err != nil {
		log.Fatalf("failed to convert markdown to HTML: %v", err)
	}
	html := buf.String()

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if _, err = io.WriteString(w, `<div class="markdown">`); err != nil {
			return err
		}
		if _, err = io.WriteString(w, html); err != nil {
			return err
		}
		if _, err = io.WriteString(w, `</div>`); err != nil {
			return err
		}
		return nil
	})
}
