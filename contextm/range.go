package contextm

import (
	"context"
)

type rangeKey struct{}

var _rangeKey rangeKey

// Keys that are set using WithRange.
func Keys(ctx context.Context) []interface{} {
	keys, ok := ctx.Value(_rangeKey).([]interface{})
	if !ok {
		return nil
	}
	return keys
}

// Range over the values that are set using WithRange.
func Range(ctx context.Context, f func(key, value interface{}) bool) {
	for _, k := range Keys(ctx) {
		if !f(k, ctx.Value(k)) {
			break
		}
	}
}

// WithRange returns a context with keys and values that can be iterated.
// The type of key should be unexported to prevent collisions, just as with regular context.WithValue calls.
func WithRange(ctx context.Context, key, val interface{}) context.Context {
	ctx = context.WithValue(ctx, key, val)
	return context.WithValue(ctx, _rangeKey, append(Keys(ctx), key))
}
