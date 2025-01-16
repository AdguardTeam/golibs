#!/bin/sh

# This comment is used to simplify checking local copies of the script.  Bump
# this number every time a significant change is made to this script.
#
# AdGuard-Project-Version: 13

verbose="${VERBOSE:-0}"
readonly verbose

if [ "$verbose" -gt '0' ]; then
	set -x
fi

# Set $EXIT_ON_ERROR to zero to see all errors.
if [ "${EXIT_ON_ERROR:-1}" -eq '0' ]; then
	set +e
else
	set -e
fi

set -f -u

# Source the common helpers, including not_found and run_linter.
. ./scripts/make/helper.sh

# Simple analyzers

# blocklist_imports is a simple check against unwanted packages.  The following
# packages are banned:
#
#   *  Package errors is replaced by our own package in the
#      github.com/AdguardTeam/golibs module.
#
#   *  Packages log and github.com/AdguardTeam/golibs/log are replaced by
#      stdlib's new package log/slog and AdGuard's new utilities package
#      github.com/AdguardTeam/golibs/logutil/slogutil.
#
#   *  Package github.com/prometheus/client_golang/prometheus/promauto is not
#      recommended, as it encourages reliance on global state.
#
#   *  Packages golang.org/x/exp/maps, golang.org/x/exp/slices, and
#      golang.org/x/net/context have been moved into stdlib.
#
#   *  Package io/ioutil is soft-deprecated.
#
#   *  Package reflect is often an overkill, and for deep comparisons there are
#      much better functions in module github.com/google/go-cmp.  Which is
#      already our indirect dependency and which may or may not enter the stdlib
#      at some point.
#
#      See https://github.com/golang/go/issues/45200.
#
#   *  Package sort is replaced by package slices.
#
#   *  Package unsafe is… unsafe.
#
# Currently, the only standard exception are files generated from protobuf
# schemas, which use package reflect.  If your project needs more exceptions,
# add and document them.
#
# NOTE:  Flag -H for grep is non-POSIX but all of Busybox, GNU, macOS, and
# OpenBSD support it.
blocklist_imports() {
	# TODO(e.burkov):  These are temporary exclusions for some packages.
	find . \
		-type 'f' \
		'(' -name '*.go' '!' -name '*.pb.go' ')' \
		'!' '(' \
		-path './cache/data.go' \
		-o -path './cache/list.go' \
		-o -path './errors/errors.go' \
		-o -path './hostsfile/storage.go' \
		-o -path './log/example_test.go' \
		-o -path './log/log.go' \
		-o -path './logutil/slogutil/legacy*.go' \
		-o -path './logutil/slogutil/slogutil.go' \
		-o -path './netutil/urlutil/url.go' \
		-o -path './testutil/log.go' \
		-o -path './testutil/testutil.go' \
		')' \
		-exec \
		'grep' \
		'-H' \
		'-e' '[[:space:]]"errors"$' \
		'-e' '[[:space:]]"github.com/AdguardTeam/golibs/log"$' \
		'-e' '[[:space:]]"github.com/prometheus/client_golang/prometheus/promauto"$' \
		'-e' '[[:space:]]"golang.org/x/exp/maps"$' \
		'-e' '[[:space:]]"golang.org/x/exp/slices"$' \
		'-e' '[[:space:]]"golang.org/x/net/context"$' \
		'-e' '[[:space:]]"io/ioutil"$' \
		'-e' '[[:space:]]"log"$' \
		'-e' '[[:space:]]"reflect"$' \
		'-e' '[[:space:]]"sort"$' \
		'-e' '[[:space:]]"unsafe"$' \
		'-n' \
		'{}' \
		';'
}

# method_const is a simple check against the usage of some raw strings and
# numbers where one should use named constants.
#
# NOTE:  Flag -H for grep is non-POSIX but all of Busybox, GNU, macOS, and
# OpenBSD support it.
method_const() {
	find . \
		-type 'f' \
		-name '*.go' \
		-exec \
		'grep' \
		'-H' \
		'-e' '"DELETE"' \
		'-e' '"GET"' \
		'-e' '"PATCH"' \
		'-e' '"POST"' \
		'-e' '"PUT"' \
		'-n' \
		'{}' \
		';'
}

# underscores is a simple check against Go filenames with underscores.  Add new
# build tags and OS as you go.  The main goal of this check is to discourage the
# use of filenames like client_manager.go.
underscores() {
	underscore_files="$(
		find . \
			-type 'f' \
			-name '*_*.go' \
			'!' '(' -name '*_darwin.go' \
			-o -name '*_generate.go' \
			-o -name '*_linux.go' \
			-o -name '*_others.go' \
			-o -name '*_test.go' \
			-o -name '*_unix.go' \
			-o -name '*_windows.go' \
			')' \
			-exec 'printf' '\t%s\n' '{}' ';'
	)"
	readonly underscore_files

	if [ "$underscore_files" != '' ]; then
		printf \
			'found file names with underscores:\n%s\n' \
			"$underscore_files"
	fi
}

# TODO(a.garipov): Add an analyzer to look for `fallthrough`, `goto`, and `new`?

# Checks

run_linter -e blocklist_imports

run_linter -e method_const

run_linter -e underscores

run_linter -e gofumpt --extra -e -l .

run_linter "${GO:-go}" vet ./...

run_linter govulncheck ./...

run_linter gocyclo --over 10 .

run_linter gocognit --over 10 .

run_linter ineffassign ./...

run_linter unparam ./...

find . \
	-type 'f' \
	'(' \
	-name 'Makefile' \
	-o -name '*.conf' \
	-o -name '*.go' \
	-o -name '*.mod' \
	-o -name '*.sh' \
	-o -name '*.yaml' \
	-o -name '*.yml' \
	')' \
	-exec 'misspell' '--error' '{}' '+'

run_linter nilness ./...

# TODO(a.garipov):  Remove or replace cache package.
run_linter fieldalignment \
	./container/ \
	./contextutil/ \
	./errors/ \
	./hostsfile/ \
	./httphdr/ \
	./internal/... \
	./ioutil/ \
	./log/ \
	./logutil/... \
	./mathutil/ \
	./netutil/... \
	./osutil/ \
	./service/ \
	./stringutil/ \
	./syncutil/ \
	./testutil/... \
	./timeutil/ \
	;

run_linter -e shadow --strict ./...

# TODO(a.garipov):  Remove or replace cache package.
run_linter gosec --exclude-dir='cache' --quiet ./...

run_linter errcheck ./...

staticcheck_matrix='
darwin:  GOOS=darwin
freebsd: GOOS=freebsd
linux:   GOOS=linux
openbsd: GOOS=openbsd
windows: GOOS=windows
'
readonly staticcheck_matrix

printf '%s' "$staticcheck_matrix" | run_linter staticcheck --matrix ./...
