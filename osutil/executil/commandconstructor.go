package executil

import (
	"context"
	"io"
	"os/exec"

	"github.com/AdguardTeam/golibs/validate"
)

// CommandConfig is the configuration for a command to run.
type CommandConfig struct {
	// Stderr, if not nil, is the process's standard error.
	//
	// See [exec.Cmd.Stderr].
	Stderr io.Writer

	// Stdin, if not nil, is the process's standard input.
	//
	// See [exec.Cmd.Stdin].
	Stdin io.Reader

	// Stdout, if not nil, is the process's standard output.
	//
	// See [exec.Cmd.Stdout].
	Stdout io.Writer

	// Path is the path to the command to run.  It should not be empty.
	Path string

	// Args, if not empty, are the command-line arguments, excluding the command
	// name itself.
	Args []string
}

// CommandConstructor is the interface for creating OS commands.
type CommandConstructor interface {
	// New creates an OS command ready to run.  conf should not be nil and
	// should be valid.  If err is not nil, c must be nil, and vice versa.
	// Implementations must document if they're using the provided context for
	// cancellation.
	New(ctx context.Context, conf *CommandConfig) (c Command, err error)
}

// EmptyCommandConstructor is a [CommandConstructor] that creates instances of
// [EmptyCommand].
type EmptyCommandConstructor struct{}

// type check
var _ CommandConstructor = EmptyCommandConstructor{}

// New implements the [CommandConstructor] interface for
// EmptyCommandConstructor.  c is always of type [EmptyCommand] and err is
// always nil.
func (EmptyCommandConstructor) New(_ context.Context, _ *CommandConfig) (c Command, err error) {
	return EmptyCommand{}, nil
}

// SystemCommandConstructor is a [CommandConstructor] that creates instances of
// [SystemCommand].
type SystemCommandConstructor struct{}

// type check
var _ CommandConstructor = SystemCommandConstructor{}

// New implements the [CommandConstructor] interface for
// SystemCommandConstructor.
//
// See [exec.CommandContext] for the documentation about the handling of ctx.
func (SystemCommandConstructor) New(
	ctx context.Context,
	conf *CommandConfig,
) (c Command, err error) {
	err = validate.NotNil("conf", conf)
	if err != nil {
		// Don't wrap the error, because it's informative enough as is.
		return nil, err
	}

	err = validate.NotEmpty("conf.Path", conf.Path)
	if err != nil {
		// Don't wrap the error, because it's informative enough as is.
		return nil, err
	}

	// #nosec G204 -- This is library code.
	execCmd := exec.CommandContext(ctx, conf.Path, conf.Args...)

	execCmd.Stderr = conf.Stderr
	execCmd.Stdin = conf.Stdin
	execCmd.Stdout = conf.Stdout

	return &SystemCommand{
		cmd: execCmd,
	}, nil
}
