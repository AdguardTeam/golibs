package service

import (
	"cmp"
	"context"
	"log/slog"
	"os"
	"slices"
	"time"

	"github.com/AdguardTeam/golibs/logutil/slogutil"
	"github.com/AdguardTeam/golibs/osutil"
)

// SignalHandler processes incoming signals and shuts services down.
type SignalHandler struct {
	logger          *slog.Logger
	signal          chan os.Signal
	refreshers      []Refresher
	services        []Interface
	refreshTimeout  time.Duration
	shutdownTimeout time.Duration
}

// SignalHandlerConfig contains the configuration for a signal handler.  See
// [NewSignalHandler].
type SignalHandlerConfig struct {
	// SignalNotifier is used to notify the handler about signals.
	//
	// If nil, [osutil.DefaultSignalNotifier] is used.
	SignalNotifier osutil.SignalNotifier

	// Logger is used for logging the handling of signals.  It should
	// include a prefix; the recommended prefix is [SignalHandlerPrefix].
	//
	// If nil, [slog.Default] with [SignalHandlerPrefix] is used.
	Logger *slog.Logger

	// RefreshTimeout is the timeout used to refresh all added refreshers.
	//
	// If zero, [SignalHandlerRefreshTimeout] is used.
	RefreshTimeout time.Duration

	// ShutdownTimeout is the timeout used to shut down all added services
	// gracefully.
	//
	// If zero, [SignalHandlerShutdownTimeout] is used.
	ShutdownTimeout time.Duration
}

// defaultSignalHandlerConf is the default configuration for a signal handler.
var defaultSignalHandlerConf = &SignalHandlerConfig{
	SignalNotifier:  osutil.DefaultSignalNotifier{},
	Logger:          slog.Default().With(slogutil.KeyPrefix, SignalHandlerPrefix),
	RefreshTimeout:  SignalHandlerRefreshTimeout,
	ShutdownTimeout: SignalHandlerShutdownTimeout,
}

// SignalHandlerPrefix is the default and recommended prefix for the logger of a
// [SignalHandler].
const SignalHandlerPrefix = "sighdlr"

// SignalHandlerRefreshTimeout is the default refresh timeout for all refreshers
// in a [SignalHandler].
const SignalHandlerRefreshTimeout = 10 * time.Second

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
		refreshers:      nil,
		services:        nil,
		refreshTimeout:  cmp.Or(c.RefreshTimeout, defaultSignalHandlerConf.RefreshTimeout),
		shutdownTimeout: cmp.Or(c.ShutdownTimeout, defaultSignalHandlerConf.ShutdownTimeout),
	}

	notifier := cmp.Or(c.SignalNotifier, defaultSignalHandlerConf.SignalNotifier)
	osutil.NotifyShutdownSignal(notifier, h.signal)
	osutil.NotifyReconfigureSignal(notifier, h.signal)

	return h
}

// AddRefresher adds refreshers to the signal handler.
//
// It must not be called concurrently with [Handle].
func (h *SignalHandler) AddRefresher(refrs ...Refresher) {
	h.refreshers = append(h.refreshers, refrs...)
}

// AddService adds services to the signal handler.
//
// It must not be called concurrently with [Handle].
func (h *SignalHandler) AddService(svcs ...Interface) {
	h.services = append(h.services, svcs...)
}

// Add adds a services to the signal handler.
//
// It must not be called concurrently with [Handle].
//
// Deprecated: Use [SignalHandler.AddService] instead.
func (h *SignalHandler) Add(svcs ...Interface) {
	h.AddService(svcs...)
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

		if osutil.IsReconfigureSignal(sig) {
			h.reconfigure(ctx)
		} else if osutil.IsShutdownSignal(sig) {
			return h.shutdown(ctx)
		}
	}

	// Shouldn't happen, since h.signal is currently never closed.
	panic("unexpected close of h.signal")
}

// reconfigure refreshes all added refreshers.
func (h *SignalHandler) reconfigure(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, h.refreshTimeout)
	defer cancel()

	h.logger.InfoContext(ctx, "reconfiguring")

	for i, r := range slices.Backward(h.refreshers) {
		err := r.Refresh(ctx)
		if err == nil {
			continue
		}

		h.logger.ErrorContext(ctx, "refreshing", "idx", i, slogutil.KeyError, err)
	}

	h.logger.InfoContext(ctx, "reconfigured")
}

// shutdown gracefully shuts down all services.  status is
// [osutil.ExitCodeSuccess] on success and [osutil.ExitCodeFailure] on error.
func (h *SignalHandler) shutdown(ctx context.Context) (status osutil.ExitCode) {
	ctx, cancel := context.WithTimeout(ctx, h.shutdownTimeout)
	defer cancel()

	h.logger.InfoContext(ctx, "shutting down")

	status = osutil.ExitCodeSuccess
	for i, s := range slices.Backward(h.services) {
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
