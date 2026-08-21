package templates

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"io"

	"github.com/a-h/templ"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"

	"monks.co/piano-alone/assets"
)

// The page inlines its stylesheet and scripts rather than linking
// them: there is one page, it is loaded once, and a round trip on a
// phone at a venue costs more than the bytes do.
//
// Each of these used to read from ./templates at process start and
// log.Fatalf on a miss, which made the working directory part of the
// deployment and a missing build artifact into a crash. They read from
// the embedded assets now, and a miss is an error on the render — the
// page that wanted the file fails, and the process keeps serving.

func style(filename string) templ.Component {
	return wrapped(filename, "<style>\n", "\n</style>\n")
}

func script(filename string) templ.Component {
	return wrapped(filename, "<script>\n", "\n</script>\n")
}

func wrapped(filename, open, close string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		input, err := assets.Read(filename)
		if err != nil {
			return fmt.Errorf("reading asset %s: %w", filename, err)
		}
		if _, err := io.WriteString(w, open); err != nil {
			return err
		}
		if _, err := w.Write(input); err != nil {
			return err
		}
		_, err = io.WriteString(w, close)
		return err
	})
}

// markdownFiles are the prose pages: the piece's own text and the
// operator's instructions. They are source rather than build output,
// so they are embedded by name.
//
//go:embed copy.md download.md
var markdownFiles embed.FS

func markdown(filename string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		input, err := markdownFiles.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("reading markdown %s: %w", filename, err)
		}
		md := goldmark.New(
			goldmark.WithRendererOptions(html.WithUnsafe()),
			goldmark.WithExtensions(extension.Typographer, extension.Linkify),
		)
		var buf bytes.Buffer
		if err := md.Convert(input, &buf); err != nil {
			return fmt.Errorf("rendering markdown %s: %w", filename, err)
		}
		if _, err := io.WriteString(w, `<div class="markdown">`); err != nil {
			return err
		}
		if _, err := w.Write(buf.Bytes()); err != nil {
			return err
		}
		_, err = io.WriteString(w, `</div>`)
		return err
	})
}
