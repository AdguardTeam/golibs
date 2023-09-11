//go:build !windows

package hostsfile

// defaultHostsPaths returns default paths to hosts files for UNIX.
func defaultHostsPaths() (paths []string, err error) {
	return []string{"etc/hosts"}, nil
}
