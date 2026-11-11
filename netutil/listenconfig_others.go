//go:build !linux

package netutil

import "github.com/AdguardTeam/golibs/errors"

// newListenConfigWithOOB returns new listen config which saves local address
// from out-of-band data.  It returns [errors.ErrUnsupported] on unsupported
// platforms.
//
// TODO(f.setrakov): Support on macOS.
func newListenConfigWithOOB(_ *ControlConfig) (lc ListenConfig, err error) {
	return nil, errors.ErrUnsupported
}
