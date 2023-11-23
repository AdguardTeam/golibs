//go:build windows

package osutil

import (
	"io/fs"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows"
)

// rootDirFS returns a filesystem rooted at the root directory of the system
// volume.  Note that it uses the C: volume as a fallback.
func rootDirFS() (fsys fs.FS) {
	// TODO(a.garipov): Use a better way if golang/go#44279 is ever resolved.
	sysDir, err := windows.GetSystemDirectory()
	if err != nil {
		// Assume that C: is the safe default.
		return os.DirFS("C:")
	}

	return os.DirFS(filepath.VolumeName(sysDir))
}
