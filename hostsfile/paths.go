package hostsfile

import "io/fs"

// justify import
var _ fs.FS

// DefaultHostsPaths returns the slice of default paths to system hosts files.
// Those are relative to the corresponding operating system's root directory and
// always use slashes to separate path elements, since those are assumed to be
// used with [fs.FS].  It may return an error only on Windows.
func DefaultHostsPaths() (p []string, err error) {
	return defaultHostsPaths()
}
