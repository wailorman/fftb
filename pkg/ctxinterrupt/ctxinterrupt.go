package ctxinterrupt

import (
	"context"
	"os"
	"os/signal"
)

func ContextWithInterruptHandling(inputCtx context.Context) context.Context {
	ctx, cancel := context.WithCancel(inputCtx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}
	}()

	return ctx
}
