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
	// Don't leak temp files left by failed write attempts
	defer func() {
		_ = tmpFile.Close()
		if err != nil {
			_ = os.Remove(tmpPath)
		}
	}()

	_, err = tmpFile.Write(data)
	if err != nil {
		return err
	}

	// Change file mode to retain compat with old version of function
	if runtime.GOOS != "windows" {
		_ = tmpFile.Chmod(0644)
	}

	// Close explicitly before renaming (this is necessary on Windows)
	_ = tmpFile.Close()

	// Assign err explicitly to make defer func aware about error
	err = os.Rename(tmpPath, path)
	return err
}
