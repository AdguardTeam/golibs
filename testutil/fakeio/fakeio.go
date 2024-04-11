// Package fakeio contains fake implementations of interfaces from package io
// from the standard library.
//
// It is recommended to fill all methods that shouldn't be called with:
//
//	panic("not implemented")
//
// in the body of the test, so that if the method is called the panic backtrace
// points to the method definition in the test.  See the package example.
package fakeio

import "io"

// Closer is the [io.Closer] implementation for tests.
type Closer struct {
	OnClose func() (err error)
}

// type check
var _ io.Closer = (*Closer)(nil)

// Close implements the [io.Closer] interface for *Closer.
func (w *Closer) Close() (err error) {
	return w.OnClose()
}

// Reader is the [io.Reader] implementation for tests.
type Reader struct {
	OnRead func(b []byte) (n int, err error)
}

// type check
var _ io.Reader = (*Reader)(nil)

// Read implements the [io.Reader] interface for *Reader.
func (w *Reader) Read(b []byte) (n int, err error) {
	return w.OnRead(b)
}

// Writer is the [io.Writer] implementation for tests.
type Writer struct {
	OnWrite func(b []byte) (n int, err error)
}

// type check
var _ io.Writer = (*Writer)(nil)

// Write implements the [io.Writer] interface for *Writer.
func (w *Writer) Write(b []byte) (n int, err error) {
	return w.OnWrite(b)
}
