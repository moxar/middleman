package httpm

import (
	"context"
	"net/http"
)

// ClientMiddleware is a function that chains caller.
type ClientMiddleware = func(Caller) Caller

// Caller is the client side interface responsible for doing HTTP calls.
type Caller interface {
	Call(ctx context.Context, method, url string, h http.Header, in, out interface{}) error
}

// CallerFunc is a function that can be used as a caller.
type CallerFunc func(ctx context.Context, method, url string, h http.Header, in, out interface{}) error

// Call implements the Caller interface.
func (c CallerFunc) Call(ctx context.Context, method, url string, h http.Header, in, out interface{}) error {
	return c(ctx, method, url, h, in, out)
}

// Chain the middlewares around the caller.
// Chain(a, b, c).Call(ctx, method, url, h, in, out) is the equivalent to (c(b(a(ctx, method, url, h, in, out)))
func Chain(caller Caller, ms ...ClientMiddleware) Caller {
	for i := len(ms) - 1; i >= 0; i-- {
		caller = ms[i](caller)
	}
	return caller
}
