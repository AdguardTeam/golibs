//go:build windows

package netutil

import (
	"syscall"

	"github.com/AdguardTeam/golibs/errors"
)

// defaultListenControlWithSO is a default [ControlFunc].  It returns
// [errors.ErrUnsupported] on Windows.
func defaultListenControlWithSO(_ syscall.RawConn) (err error) {
	return errors.ErrUnsupported
}
