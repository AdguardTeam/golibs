#!/bin/sh

# This comment is used to simplify checking local copies of the script.  Bump
# this number every time a significant change is made to this script.
#
# AdGuard-Project-Version: 5

verbose="${VERBOSE:-0}"
readonly verbose

# Verbosity levels:
#   0 = Don't print anything except for errors.
#   1 = Print commands, but not nested commands.
#   2 = Print everything.
if [ "$verbose" -gt '1' ]; then
	set -x
	v_flags='-v=1'
	x_flags='-x=1'
elif [ "$verbose" -gt '0' ]; then
	set -x
	v_flags='-v=1'
	x_flags='-x=0'
else
	set +x
	v_flags='-v=0'
	x_flags='-x=0'
fi
readonly v_flags x_flags

set -e -f -u

if [ "${RACE:-1}" -eq '0' ]; then
	race_flags='--race=0'
else
	race_flags='--race=1'
fi
readonly race_flags

count_flags='--count=2'
cover_flags='--coverprofile=./cover.out'
go="${GO:-go}"
shuffle_flags='--shuffle=on'
timeout_flags="${TIMEOUT_FLAGS:---timeout=30s}"
readonly count_flags cover_flags go shuffle_flags timeout_flags

"$go" test \
	"$count_flags" \
	"$cover_flags" \
	"$race_flags" \
	"$shuffle_flags" \
	"$timeout_flags" \
	"$v_flags" \
	"$x_flags" \
	./...
