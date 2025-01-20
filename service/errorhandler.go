package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/logutil/slogutil"
	"github.com/AdguardTeam/golibs/validate"
)

// ErrorHandler is the interface for entities that handle errors from refreshes
// and other asyncrhonous jobs.
//
// TODO(a.garipov):  Consider moving to package errors.
type ErrorHandler interface {
	// Handle handles the error.  err must not be nil.
	Handle(ctx context.Context, err error)
}

// ErrorHandlerFunc is an adapter to allow the use of ordinary functions as
// [ErrorHandler]s.
type ErrorHandlerFunc func(ctx context.Context, err error)

// type check
var _ ErrorHandler = ErrorHandlerFunc(nil)

// Handle implements the [ErrorHandler] interface for ErrorHandlerFunc.
func (f ErrorHandlerFunc) Handle(ctx context.Context, err error) {
	f(ctx, err)
}

// IgnoreErrorHandler is a [ErrorHandler] that ignores all errors.
type IgnoreErrorHandler struct{}

// type check
var _ ErrorHandler = IgnoreErrorHandler{}

// Handle implements the [ErrorHandler] interface for IgnoreErrorHandler.
func (IgnoreErrorHandler) Handle(_ context.Context, _ error) {}

// SlogErrorHandler logs errors.
type SlogErrorHandler struct {
	logger *slog.Logger
	lvl    slog.Leveler
	msg    string
}

// NewSlogErrorHandler returns a new logging error handler.  l and lvl must not
// be nil.
func NewSlogErrorHandler(l *slog.Logger, lvl slog.Leveler, msg string) (h *SlogErrorHandler) {
	err := errors.Join(
		validate.NotNil("l", l),
		validate.NotNilInterface("lvl", lvl),
	)
	if err != nil {
		panic(fmt.Errorf("service.NewSlogErrorHandler: %w", err))
	}

	return &SlogErrorHandler{
		logger: l,
		lvl:    lvl,
		msg:    msg,
	}
}

// type check
var _ ErrorHandler = (*SlogErrorHandler)(nil)

// Handle implements the [ErrorHandler] interface for *SlogErrorHandler.
func (h *SlogErrorHandler) Handle(ctx context.Context, err error) {
	h.logger.Log(ctx, h.lvl.Level(), h.msg, slogutil.KeyError, err)
}
