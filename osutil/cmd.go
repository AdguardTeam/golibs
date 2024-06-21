package osutil

// ExitCode is a semantic alias for int when it's used as an exit code.
type ExitCode = int

// Exit status constants.
const (
	ExitCodeSuccess       ExitCode = 0
	ExitCodeFailure       ExitCode = 1
	ExitCodeArgumentError ExitCode = 2
)
