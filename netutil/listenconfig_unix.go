//go:build unix

package netutil

import (
	"syscall"

	"github.com/AdguardTeam/golibs/errors"
	"golang.org/x/sys/unix"
)

// defaultListenControlWithSO is a default [ControlFunc] that sets
// SO_REUSEPORT=1 only.  c must not be nil.
func defaultListenControlWithSO(c syscall.RawConn) (err error) {
	var opErr error
	err = c.Control(func(fd uintptr) {
		fdInt := int(fd)
		opErr = unix.SetsockoptInt(fdInt, unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
	})
	err = errors.Annotate(err, "setting SO_REUSEPORT: %w")

	return errors.WithDeferred(opErr, err)
}
