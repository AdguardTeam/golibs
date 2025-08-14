// Package fakeexec contains fake implementations of interfaces from package
// osutil/executil.
//
// It is recommended to fill all methods that shouldn't be called with:
//
//	panic(testutil.UnexpectedCall(arg1, arg2))
package fakeexec

import (
	"context"

	"github.com/AdguardTeam/golibs/osutil/executil"
	"github.com/AdguardTeam/golibs/testutil"
)

// Command is the [executil.Command] for tests.
type Command struct {
	OnCancel func(ctx context.Context) (err error)
	OnStart  func(ctx context.Context) (err error)
	OnWait   func(ctx context.Context) (err error)
}

// type check
var _ executil.Command = (*Command)(nil)

// Cancel implements the [executil.Command] interface for *Command.
func (c *Command) Cancel(ctx context.Context) (err error) {
	return c.OnCancel(ctx)
}

// Start implements the [executil.Command] interface for *Command.
func (c *Command) Start(ctx context.Context) (err error) {
	return c.OnStart(ctx)
}

// Wait implements the [executil.Command] interface for *Command.
func (c *Command) Wait(ctx context.Context) (err error) {
	return c.OnWait(ctx)
}

// NewCommand returns a new *Command all methods of which panic.
func NewCommand() (c *Command) {
	return &Command{
		OnCancel: func(ctx context.Context) (err error) { panic(testutil.UnexpectedCall(ctx)) },
		OnStart:  func(ctx context.Context) (err error) { panic(testutil.UnexpectedCall(ctx)) },
		OnWait:   func(ctx context.Context) (err error) { panic(testutil.UnexpectedCall(ctx)) },
	}
}

// CommandConstructor is the [executil.CommandConstructor] for tests.
type CommandConstructor struct {
	OnNew func(ctx context.Context, conf *executil.CommandConfig) (c executil.Command, err error)
}

// type check
var _ executil.CommandConstructor = (*CommandConstructor)(nil)

// New implements the [executil.CommandConstructor] interface for
// *CommandConstructor.
func (cons *CommandConstructor) New(
	ctx context.Context,
	conf *executil.CommandConfig,
) (c executil.Command, err error) {
	return cons.OnNew(ctx, conf)
}
