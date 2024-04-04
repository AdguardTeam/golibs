//go:build unix

package service

import (
	"cmp"
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/AdguardTeam/golibs/logutil/slogutil"
	"github.com/AdguardTeam/golibs/osutil"
	"golang.org/x/sys/unix"
)

// SignalHandler processes incoming signals and shuts services down.
//
// TODO(a.garipov): Expand to Windows.
type SignalHandler struct {
	logger          *slog.Logger
	signal          chan os.Signal
	services        []Interface
	shutdownTimeout time.Duration
}

// SignalHandlerConfig contains the configuration for a signal handler.  See
// [NewSignalHandler].
type SignalHandlerConfig struct {
	// SignalNotifier is used to notify the handler about signals.
	//
	// If nil, [osutil.DefaultSignalNotifier] is used.
	SignalNotifier osutil.SignalNotifier

	// Logger is used for logging the shutting down of services.  It should
	// include a prefix; the recommended prefix is [SignalHandlerPrefix].
	//
	// If nil, [slog.Default] with [SignalHandlerPrefix] is used.
	Logger *slog.Logger

	// ShutdownTimeout is the timeout used to shut down all services gracefully.
	//
	// If zero, [DefaultShutdownTimeout] is used.
	ShutdownTimeout time.Duration
}

// defaultSignalHandlerConf is the default configuration for a signal handler.
var defaultSignalHandlerConf = &SignalHandlerConfig{
	SignalNotifier:  osutil.DefaultSignalNotifier{},
	Logger:          slog.Default().With(slogutil.KeyPrefix, SignalHandlerPrefix),
	ShutdownTimeout: SignalHandlerShutdownTimeout,
}

// SignalHandlerPrefix is the default and recommended prefix for the logger of a
// [SignalHandler].
const SignalHandlerPrefix = "sighdlr"

// SignalHandlerShutdownTimeout is the default shutdown timeout for all services
// in a [SignalHandler].
const SignalHandlerShutdownTimeout = 10 * time.Second

// NewSignalHandler returns a new properly initialized *SignalHandler that shuts
// down services.  If c is nil, the defaults of [SignalHandlerConfig] are used.
func NewSignalHandler(c *SignalHandlerConfig) (h *SignalHandler) {
	if c == nil {
		c = defaultSignalHandlerConf
	}

	h = &SignalHandler{
		logger:          cmp.Or(c.Logger, defaultSignalHandlerConf.Logger),
		signal:          make(chan os.Signal, 1),
		services:        nil,
		shutdownTimeout: cmp.Or(c.ShutdownTimeout, defaultSignalHandlerConf.ShutdownTimeout),
	}

	// TODO(a.garipov): Expand these to Windows.
	notifier := cmp.Or(c.SignalNotifier, defaultSignalHandlerConf.SignalNotifier)
	notifier.Notify(h.signal, unix.SIGINT, unix.SIGQUIT, unix.SIGTERM)

	return h
}

// Add adds a services to the signal handler.
//
// It must not be called concurrently with [Handle].
func (h *SignalHandler) Add(svcs ...Interface) {
	h.services = append(h.services, svcs...)
}

// Handle processes signals from the handler's [osutil.SignalNotifier].  It
// blocks until a termination signal is received, after which it shuts down all
// services.  ctx is used for logging and serves as the base for the shutdown
// timeout.  status is [osutil.ExitCodeSuccess] on success and
// [osutil.ExitCodeFailure] on error.
//
// Handle must not be called concurrently with [Add].
func (h *SignalHandler) Handle(ctx context.Context) (status osutil.ExitCode) {
	defer slogutil.RecoverAndLog(ctx, h.logger)

	for sig := range h.signal {
		h.logger.InfoContext(ctx, "received", "signal", sig)

		switch sig {
		case
			unix.SIGINT,
			unix.SIGQUIT,
			unix.SIGTERM:

			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, h.shutdownTimeout)
			defer cancel()

			return h.shutdown(ctx)
		}
	}

	// Shouldn't happen, since h.signal is currently never closed.
	panic("unexpected close of h.signal")
}

// shutdown gracefully shuts down all services.  status is
// [osutil.ExitCodeSuccess] on success and [osutil.ExitCodeFailure] on error.
func (h *SignalHandler) shutdown(ctx context.Context) (status osutil.ExitCode) {
	h.logger.InfoContext(ctx, "shutting down")

	status = osutil.ExitCodeSuccess
	for i := len(h.services) - 1; i >= 0; i-- {
		s := h.services[i]
		err := s.Shutdown(ctx)
		if err == nil {
			continue
		}

		h.logger.ErrorContext(ctx, "shutting down service", "idx", i, slogutil.KeyError, err)

		status = osutil.ExitCodeFailure
	}

	h.logger.InfoContext(ctx, "shut down", "status", status)

	return status
}
