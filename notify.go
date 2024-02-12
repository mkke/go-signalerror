package signalerror

import (
	"context"
	"os"
	"os/signal"
)

// NotifyContext returns a copy of the parent context that is marked done
// (its Done channel is closed) when one of the listed signals arrives,
// when the returned stop function is called, or when the parent context's
// Done channel is closed, whichever happens first.
//
// It mirrors signal.NotifyContext's functionality, but additionally sets a SignalError
// as context cancel cause.
func NotifyContext(parent context.Context, signals ...os.Signal) (ctx context.Context, stop context.CancelFunc) {
	ctx, cancel := context.WithCancelCause(parent)
	c := &signalCtx{
		Context: ctx,
		cancel:  cancel,
		signals: signals,
	}
	c.ch = make(chan os.Signal, 1)
	signal.Notify(c.ch, c.signals...)
	if ctx.Err() == nil {
		go func() {
			select {
			case sig := <-c.ch:
				c.cancel(NewSignalError(sig))
			case <-c.Done():
			}
		}()
	}
	return c, c.stop
}

type signalCtx struct {
	context.Context

	cancel  context.CancelCauseFunc
	signals []os.Signal
	ch      chan os.Signal
}

func (c *signalCtx) stop() {
	c.cancel(nil)
	signal.Stop(c.ch)
}
