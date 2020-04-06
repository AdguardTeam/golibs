// Package file provides helper functions for working with files
package file

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

// SafeWrite writes data to a temporary file and then renames it to what's specified in path
func SafeWrite(path string, data []byte) error {
	dir := filepath.Dir(path)

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	// Ensure multiple simultaneous callers will not choose the same file.
	tmpFile, err := ioutil.TempFile(dir, "tmp")
	if err != nil {
		return err
	}

	tmpPath := tmpFile.Name()

	_, err = tmpFile.Write(data)
	if err != nil {
        // Don't leak temp files left by failed write attempts
		_ = tmpFile.Close()
        _ = os.Remove(tmpPath)
		return err
	}

	// Change file mode to retain compat with old version of function
	if runtime.GOOS != "windows" {
		_ = tmpFile.Chmod(0644)
	}

    _ = tmpFile.Close()
	err = os.Rename(tmpPath, path)
    if err != nil {
        _ = os.Remove(tmpPath)
    }
	return err
}
