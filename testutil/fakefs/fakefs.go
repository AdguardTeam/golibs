// Package fakefs contains fake implementations of interfaces from package io/fs
// from the standard library.
//
// It is recommended to fill all methods that shouldn't be called with:
//
//	panic("not implemented")
//
// in the body of the test, so that if the method is called the panic backtrace
// points to the method definition in the test.  See the package example.
package fakefs

import "io/fs"

// File is the [fs.File] for tests.
type File struct {
	OnClose func() error
	OnRead  func(b []byte) (n int, err error)
	OnStat  func() (fi fs.FileInfo, err error)
}

// type check
var _ fs.File = (*File)(nil)

// Close implements the [fs.File] interface for *File.
func (f *File) Close() (err error) {
	return f.OnClose()
}

// Read implements the [fs.File] interface for *File.
func (f *File) Read(b []byte) (n int, err error) {
	return f.OnRead(b)
}

// Stat implements the [fs.File] interface for *File.
func (f *File) Stat() (fi fs.FileInfo, err error) {
	return f.OnStat()
}

// FS is the [fs.FS] for tests.
type FS struct {
	OnOpen func(name string) (fs.File, error)
}

// type check
var _ fs.FS = (*FS)(nil)

// Open implements the [fs.FS] interface for *FS.
func (fsys *FS) Open(name string) (f fs.File, err error) {
	return fsys.OnOpen(name)
}

// type check
var _ fs.GlobFS = (*GlobFS)(nil)

// GlobFS is the [fs.GlobFS] for tests.
type GlobFS struct {
	OnOpen func(name string) (f fs.File, err error)
	OnGlob func(pattern string) (paths []string, err error)
}

// Open implements the [fs.GlobFS] interface for *GlobFS.
func (fsys *GlobFS) Open(name string) (f fs.File, err error) {
	return fsys.OnOpen(name)
}

// Glob implements the [fs.GlobFS] interface for *GlobFS.
func (fsys *GlobFS) Glob(pattern string) (paths []string, err error) {
	return fsys.OnGlob(pattern)
}

// type check
var _ fs.StatFS = (*StatFS)(nil)

// StatFS is the [fs.StatFS] for tests.
type StatFS struct {
	OnOpen func(name string) (f fs.File, err error)
	OnStat func(name string) (fi fs.FileInfo, err error)
}

// Open implements the [fs.StatFS] interface for *StatFS.
func (fsys *StatFS) Open(name string) (f fs.File, err error) {
	return fsys.OnOpen(name)
}

// Stat implements the [fs.StatFS] interface for *StatFS.
func (fsys *StatFS) Stat(name string) (fi fs.FileInfo, err error) {
	return fsys.OnStat(name)
}
