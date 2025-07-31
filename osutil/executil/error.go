package executil

import (
	"os/exec"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/osutil"
)

// ExitCodeError is the interface for errors returned by methods of [Command]
// that have an associated exit code.
type ExitCodeError interface {
	error

	// ExitCode must return the exit code of the exited process, or -1 if the
	// process hasn't exited or was terminated by a signal.
	//
	// See [os.ProcessState.ExitCode].
	ExitCode() (c osutil.ExitCode)
}

// Make sure that *exec.ExitError implements the interface.
var _ ExitCodeError = (*exec.ExitError)(nil)

// ExitCodeFromError returns an exit code and true if err contains an
// [ExitCodeError].  Otherwise, it returns 0 and false.
func ExitCodeFromError(err error) (c osutil.ExitCode, ok bool) {
	var exitErr ExitCodeError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode(), true
	}

	return 0, false
}
