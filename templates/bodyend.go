package templates

import (
	"context"
	"html/template"
	"io"

	"github.com/a-h/templ"
)

// The pages carry their own design — no shared chrome, no design
// system — but a host still needs somewhere to put the markup its own
// pages all carry. BodyEnd is that somewhere: whatever it returns is
// written just before </body> on every page.
//
// It travels in the context rather than as a parameter because it
// belongs to the request, not to the page: every component here would
// otherwise grow an argument it does not use, to hand to a layout it
// does not know about.

type bodyEndKey struct{}

// BodyEnd renders host-owned markup at the end of a page's body. It
// takes no framework types, so any code that can produce HTML can
// supply one.
type BodyEnd func() template.HTML

// WithBodyEnd returns a context whose pages end with f's markup. A
// context without one renders nothing extra, which is what the
// standalone server does.
func WithBodyEnd(ctx context.Context, f BodyEnd) context.Context {
	return context.WithValue(ctx, bodyEndKey{}, f)
}

func bodyEnd() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		f, ok := ctx.Value(bodyEndKey{}).(BodyEnd)
		if !ok || f == nil {
			return nil
		}
		_, err := io.WriteString(w, string(f()))
		return err
	})
}
