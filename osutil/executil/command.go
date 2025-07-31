package executil

import (
	"context"
	"os/exec"
)

// Command is the interface for OS commands.
//
// Methods must return an error containing an [ExitCodeError] if the error
// signals that the process has exited with a non-zero exit code.
type Command interface {
	// Cancel must stop the process and make all Wait calls exit.  It must not
	// be called before the call to Start.
	//
	// Implementations must document if they're using the provided context for
	// cancellation or not.
	Cancel(ctx context.Context) (err error)

	// Start must start the execution of the command but not wait for it to
	// complete.  Start must only be called once.  After a successful call to
	// Start, Cancel or Wait must be called.
	//
	// Implementations must document if they're using the provided context for
	// cancellation or not.
	Start(ctx context.Context) (err error)

	// Wait waits for the command to exit.  Wait must only be called after
	// Start.
	//
	// Implementations must document if they're using the provided context for
	// cancellation or not.
	Wait(ctx context.Context) (err error)
}

// EmptyCommand is a [Command] that does nothing.  Its methods ignore the
// context and the returned errors are always nil.
type EmptyCommand struct{}

// type check
var _ Command = EmptyCommand{}

// Cancel implements the [Command] interface for EmptyCommand.
func (EmptyCommand) Cancel(_ context.Context) (err error) { return nil }

// Start implements the [Command] interface for EmptyCommand.
func (EmptyCommand) Start(_ context.Context) (err error) { return nil }

// Wait implements the [Command] interface for EmptyCommand.
func (EmptyCommand) Wait(_ context.Context) (err error) { return nil }

// SystemCommand is a [Command] that uses [*exec.Cmd] to execute system
// commands.
//
// TODO(a.garipov):  Consider actually using the context in the methods for
// cancellation.
type SystemCommand struct {
	cmd *exec.Cmd
}

// type check
var _ Command = (*SystemCommand)(nil)

// Cancel implements the [Command] interface for *SystemCommand.  The context is
// not used for cancellation.
func (c *SystemCommand) Cancel(_ context.Context) (err error) {
	return c.cmd.Cancel()
}

// Start implements the [Command] interface for *SystemCommand.  The context is
// not used for cancellation.
func (c *SystemCommand) Start(_ context.Context) (err error) {
	return c.cmd.Start()
}

// Wait implements the [Command] interface for *SystemCommand.  The context is
// not used for cancellation.
func (c *SystemCommand) Wait(_ context.Context) (err error) {
	return c.cmd.Wait()
}
