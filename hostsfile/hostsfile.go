// Package hostsfile provides utilities for working with system hosts files.
// The syntax of the hosts files described in man page [hosts(5)], with
// hostname's syntax from [RFC-952], including its updates from [RFC-1123] and
// further ones.
//
// [hosts(5)]: https://man7.org/linux/man-pages/man5/hosts.5.html
// [RFC-952]: https://datatracker.ietf.org/doc/html/rfc952
// [RFC-1123]: https://datatracker.ietf.org/doc/html/rfc1123
package hostsfile

import (
	"context"
	"io"
)

// NamedReader is an optional interface that may be implemented by an
// [io.Reader] to provide the name of the data source.
type NamedReader interface {
	io.Reader

	// Name returns the name of the data source.
	Name() (name string)
}

// Set handles successfully unmarshaled records.
type Set interface {
	// Add adds rec to the set.  rec should be valid.
	Add(ctx context.Context, rec *Record)
}

// DiscardSet is a [Set] that discards all records.
type DiscardSet struct{}

// type check
var _ Set = (*DiscardSet)(nil)

// Add implements the [Set] interface for DiscardSet.
func (DiscardSet) Add(_ context.Context, _ *Record) {}

// FuncSet is a functional [Set] implementation.
type FuncSet func(ctx context.Context, rec *Record)

// type check
var _ Set = FuncSet(nil)

// Add implements the [Set] interface for FuncSet.
func (f FuncSet) Add(ctx context.Context, rec *Record) { f(ctx, rec) }

// HandleSet is a [Set] that handles invalid records.
type HandleSet interface {
	Set

	// HandleInvalid unmarshals invalid records according to the err returned by
	// [Record.UnmarshalText].  data is the original line from the hosts file,
	// including spaces, srcName is the name of the data source, if provided.
	HandleInvalid(ctx context.Context, srcName string, data []byte, err error)
}
