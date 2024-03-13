package sigctx

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func New() context.Context {
	ctx, _ := NewWithCancel()
	return ctx
}

func NewWithCancel() (context.Context, func(err error)) {
	ctx, cancel := context.WithCancelCause(context.Background())
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
		sig := <-sigs
		cancel(fmt.Errorf("got signal: %s", sig))
	}()
	return ctx, cancel
}

