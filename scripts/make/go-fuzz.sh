#!/bin/sh

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

if [ "$GOMAXPROCS" = '' ]; then
	export GOMAXPROCS='2'
fi

set -e -f -u

if [ "${RACE:-1}" -eq '0' ]; then
	race_flags='--race=0'
else
	race_flags='--race=1'
fi
readonly race_flags

count_flags='--count=2'
fuzztime_flags="${FUZZTIME_FLAGS:---fuzztime=10s}"
go="${GO:-go}"
shuffle_flags='--shuffle=on'
timeout_flags="${TIMEOUT_FLAGS:---timeout=30s}"
readonly count_flags fuzztime_flags go shuffle_flags timeout_flags

# TODO(a.garipov): File an issue about using --fuzz with multiple packages.
while read -r pkg fuzzname; do
	"$go" test \
		"$count_flags" \
		"$shuffle_flags" \
		"$race_flags" \
		"$timeout_flags" \
		"$x_flags" \
		"$v_flags" \
		"$fuzztime_flags" \
		"--fuzz=${fuzzname}" \
		"$pkg" \
		;
done <<-'EOF'
	./timeutil  FuzzDuration_МarshalText_roundTrip
	./timeutil  FuzzDuration_UnmarshalText
	./netutil   FuzzIsValidHostname$
	./netutil   FuzzIsValidHostnameLabel
	./netutil   FuzzIsValidIPPortString
	./netutil   FuzzIsValidIPPrefixString
	./netutil   FuzzIsValidIPString
	./netutil   FuzzIsValidMACString
	./hostsfile FuzzParse
	./hostsfile FuzzRecord_UnmarshalText
	./netutil   FuzzSubnetSet_Contains_v4
	./netutil   FuzzSubnetSet_Contains_v6
EOF
