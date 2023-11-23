package osutil

import "io/fs"

// RootDirFS returns a filesystem rooted at the system's root directory.
func RootDirFS() (fsys fs.FS) {
	return rootDirFS()
}
