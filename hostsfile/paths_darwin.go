//go:build darwin

package hostsfile

// defaultHostsPaths returns default paths to hosts files for macOS.
func defaultHostsPaths() (paths []string, err error) {
	return []string{"private/etc/hosts"}, nil
}
