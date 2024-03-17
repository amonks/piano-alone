package templates

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/a-h/templ"
)

func script(filename string) templ.Component {
	input, err := os.ReadFile(filepath.Join("templates", filename))
	if err != nil {
		log.Fatalf("could not find script file: templates/%s", filename)
	}

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if _, err := io.WriteString(w, "<script>\n"); err != nil {
			return err
		}
		if _, err := w.Write(input); err != nil {
			return err
		}
		if _, err := io.WriteString(w, "\n</script>\n"); err != nil {
			return err
		}
		return nil
	})
}
