package requestid

import (
	"context"
	"fmt"

	"github.com/AdguardTeam/golibs/errors"
)

// ctxKey is the type for context keys.
type ctxKey int

// Context key values.
const (
	ctxKeyRequestID ctxKey = iota
)

// ContextWithRequestID returns a new context with the given ID.
func ContextWithRequestID(parent context.Context, id ID) (ctx context.Context) {
	return context.WithValue(parent, ctxKeyRequestID, id)
}

// IDFromContext returns ID for this request, if any.
func IDFromContext(ctx context.Context) (id ID, ok bool) {
	v := ctx.Value(ctxKeyRequestID)
	if v == nil {
		return id, false
	}

	id, ok = v.(ID)
	if !ok {
		panic(fmt.Errorf("bad type for ctxKeyRequestID: %T(%[1]v)", v))
	}

	return id, true
}

// MustIDFromContext returns ID for this request and panics if there is no ID.
func MustIDFromContext(ctx context.Context) (id ID) {
	id, ok := IDFromContext(ctx)
	if !ok {
		panic(fmt.Errorf("request id in context: %w", errors.ErrNoValue))
	}

	return id
}
