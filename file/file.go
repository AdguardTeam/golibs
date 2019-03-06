// File functions

package file

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// SafeWrite writes data to a temporary file and then renames it to what's specified in path
func SafeWrite(path string, data []byte) error {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	tmpPath := path + ".tmp"
	err = ioutil.WriteFile(tmpPath, data, 0644)
	if err != nil {
		return err
	}

	return os.Rename(tmpPath, path)
}
