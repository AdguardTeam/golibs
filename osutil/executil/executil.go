// Package executil contains utilities for running OS commands.
//
// TODO(a.garipov):  Consider adding metrics and telemetry.
package executil

import (
	"bytes"
	"context"
	"fmt"

	"github.com/AdguardTeam/golibs/ioutil"
	"github.com/c2h5oh/datasize"
)

// Run is a utility function that constructs a command, starts it, and waits for
// it to complete.  cons must not be nil.  conf must not be nil and must be
// valid.
func Run(ctx context.Context, cons CommandConstructor, conf *CommandConfig) (err error) {
	c, err := cons.New(ctx, conf)
	if err != nil {
		return fmt.Errorf("constructing: %w", err)
	}

	err = c.Start(ctx)
	if err != nil {
		return fmt.Errorf("starting: %w", err)
	}

	err = c.Wait(ctx)
	if err != nil {
		return fmt.Errorf("running: %w", err)
	}

	return nil
}

// RunWithPeek is a utility function that constructs a command, starts it, and
// waits for it to complete.  If there is any error, it is wrapped with the
// first errLimit bytes of stderr and stdout.  cmdPath should not be empty.
func RunWithPeek(
	ctx context.Context,
	cons CommandConstructor,
	errLimit datasize.ByteSize,
	cmdPath string,
	args ...string,
) (err error) {
	stdoutBuf := &bytes.Buffer{}
	stderrBuf := &bytes.Buffer{}

	stderr := ioutil.NewTruncatedWriter(stderrBuf, errLimit.Bytes())
	stdout := ioutil.NewTruncatedWriter(stdoutBuf, errLimit.Bytes())

	err = Run(ctx, cons, &CommandConfig{
		Path:   cmdPath,
		Args:   args,
		Stderr: stderr,
		Stdout: stdout,
	})
	if err != nil {
		return fmt.Errorf(
			"%w; stderr peek: %q; stdout peek: %q",
			err,
			stderrBuf,
			stdoutBuf,
		)
	}

	return nil
}
