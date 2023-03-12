package ctxinterrupt

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func ContextWithInterruptHandling(inputCtx context.Context) context.Context {
	ctx, cancel := context.WithCancel(inputCtx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-c
		cancel()
	}()

	return ctx
}
