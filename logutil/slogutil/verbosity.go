package slogutil

import (
	"fmt"
	"log/slog"

	"github.com/AdguardTeam/golibs/errors"
)

// Acceptable [slog.Level] levels.
const (
	LevelTrace = slog.Level(-8)
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)

// VerbosityToLevel returns log level for given verbosity.
func VerbosityToLevel(l uint8) (lvl slog.Level, err error) {
	switch l {
	case 0:
		return LevelInfo, nil
	case 1:
		return LevelDebug, nil
	case 2:
		return LevelTrace, nil
	default:
		return lvl, fmt.Errorf("%w: %d", errors.ErrBadEnumValue, l)
	}
}
