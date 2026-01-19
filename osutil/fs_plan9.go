//go:build plan9

package osutil

import (
	"io/fs"
	"os"
)

// rootDirFS returns a filesystem rooted at the system's root directory.
func rootDirFS() (fsys fs.FS) {
	return os.DirFS("/")
}
