// Package version contains the version information.
//
// Example of the build command:
//
//	PKG="github.com/AdguardTeam/golibs/version"
//
//	go build -ldflags " \
//		-X ${PKG}.version=1.0.0 \
//		-X ${PKG}.committime=2026-04-20T10:00:00Z \
//		-X ${PKG}.revision="e3133f34fsdsd34fsfdbxffd54edf2123cds1ds0" \
//		-X ${PKG}.branch="AGDNS-0000-test" \
//		"
package version

import "os"

// These are set by the linker.  Unfortunately, we cannot set constants during
// linking, and Go doesn't have a concept of immutable variables, so to be
// thorough we have to only export them through getters.
var (
	branch     string
	committime string
	revision   string
	version    string
)

// Branch returns the compiled-in value of the Git branch.
func Branch() (b string) {
	return branch
}

// CommitTime returns the compiled-in value of the commit time as a string.
func CommitTime() (t string) {
	return committime
}

// Revision returns the compiled-in value of the Git revision.
func Revision() (r string) {
	return revision
}

// Version returns the compiled-in value of the build version as a string. It
// returns either the compiled-in value or the one from the environment.
func Version() (v string) {
	v, ok := os.LookupEnv("APP_RUNTIME_VERSION")
	if ok {
		return v
	}

	return version
}
