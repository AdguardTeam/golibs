// Package fakeexec contains fake implementations of interfaces from package
// osutil/executil.
//
// It is recommended to fill all methods that shouldn't be called with:
//
//	panic("not implemented")
//
// in the body of the test, so that if the method is called the panic backtrace
// points to the method definition in the test.
package fakeexec

import (
	"context"
	"fmt"

	"github.com/AdguardTeam/golibs/osutil/executil"
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
		OnCancel: func(_ context.Context) (err error) {
			panic(fmt.Errorf("unexpected call to fakeexec.(*Command).Cancel()"))
		},
		OnStart: func(_ context.Context) (err error) {
			panic(fmt.Errorf("unexpected call to fakeexec.(*Command).Start()"))
		},
		OnWait: func(_ context.Context) (err error) {
			panic(fmt.Errorf("unexpected call to fakeexec.(*Command).Wait()"))
		},
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
