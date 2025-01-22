// Package sentryutil contains utilities for functions for working with Sentry.
package sentryutil

import (
	"cmp"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/logutil/slogutil"
	"github.com/getsentry/sentry-go"
)

// DefaultLoggerPrefix is the default prefix for the debug logger set in
// [SetDefaultLogger].
const DefaultLoggerPrefix = "sentry_default_hub"

// DSNEnv is the name of the environment variable for the Sentry DSN.
const DSNEnv = "SENTRY_DSN"

// InitDefaultHub initializes the default Sentry hub.  It returns an error if
// [DSNEnv] isn't set in the environment.
func InitDefaultHub(release string) (err error) {
	dsn := os.Getenv(DSNEnv)
	if dsn == "" {
		return fmt.Errorf("env %s: %w", DSNEnv, errors.ErrEmptyValue)
	}

	return sentry.Init(sentry.ClientOptions{
		Dsn:     dsn,
		Release: release,
	})
}

// MustInitDefaultHub is like [InitDefaultHub] but panics on errors.  This
// function should be used in a main function.
func MustInitDefaultHub(release string) {
	errors.Check(InitDefaultHub(release))
}

// ReportPanics reports all panics using the default Sentry hub and repanics.
// This function should be used in a main function, after the default Sentry hub
// has been configured.  It should be called in a defer.
func ReportPanics() {
	v := recover()
	if v == nil {
		return
	}

	sentry.CaptureException(errors.FromRecovered(v))
	sentry.Flush(1 * time.Second)

	panic(v)
}

// SetDefaultLogger sets the default Sentry logger to l with prefix p and level
// debug.  If p is empty, [DefaultLoggerPrefix] is used.
func SetDefaultLogger(l *slog.Logger, p string) {
	p = cmp.Or(p, DefaultLoggerPrefix)
	l = l.With(slogutil.KeyPrefix, p)

	h := &traceMsgHandler{
		handler: l.Handler(),
	}

	sentry.Logger = slog.NewLogLogger(h, slog.LevelDebug)
}
