//go:build windows

package hostsfile

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows"
)

// defaultHostsPaths returns default paths to hosts files for Windows.
func defaultHostsPaths() (paths []string, err error) {
	sysDir, err := windows.GetSystemDirectory()
	if err != nil {
		return []string{}, fmt.Errorf("getting system directory: %w", err)
	}

	// Split all the elements of the path to join them afterwards.  This is
	// needed to make the Windows-specific path string returned by
	// [windows.GetSystemDirectory] to be compatible with [fs.FS].
	pathElems := strings.Split(sysDir, string(os.PathSeparator))
	if len(pathElems) > 0 && pathElems[0] == filepath.VolumeName(sysDir) {
		pathElems = pathElems[1:]
	}

	p := path.Join(append(pathElems, "drivers", "etc", "hosts")...)

	return []string{p}, nil
}
