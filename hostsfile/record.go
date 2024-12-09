package hostsfile

import (
	"bytes"
	"encoding"
	"fmt"
	"net/netip"
	"slices"
	"strings"

	"github.com/AdguardTeam/golibs/netutil"
)

// Record represents a single hosts file record.
type Record struct {
	// Addr is the IP address of the record.
	Addr netip.Addr

	// Source is the name of the source hosts file.
	Source string

	// Names are the hostnames for the Addr of the record.
	Names []string
}

// spaces are the expected space characters in the hosts file lines.  Note that
// those doesn't include newline characters, and also those are ASCII only.
const spaces = " \t"

// type check
var _ encoding.TextUnmarshaler = (*Record)(nil)

// UnmarshalText implements the [encoding.TextUnmarshaler] interface for
// *Record.  It only returns the following errors:
//   - [ErrEmptyLine] if the line is empty or contains only spaces and comments;
//   - [ErrNoHosts] if the record doesn't contain any space delimiters, but the
//     IP address may appear valid;
//   - error from [netip.ParseAddr] on invalid IP address;
//   - [netutil.AddrError] if one of the hostnames is invalid, but the Record
//     may already contain hostnames.
//
// Note that this function doesn't set the Source field of rec, see [Parse] and
// [HandleSet] for details.
func (rec *Record) UnmarshalText(data []byte) (err error) {
	if commIdx := bytes.IndexByte(data, '#'); commIdx >= 0 {
		// Trim comment.
		data = data[:commIdx]
	}

	field, data := cutField(bytes.Trim(data, spaces))
	if len(field) == 0 {
		// Empty line.
		return ErrEmptyLine
	} else if len(data) == 0 {
		// The only field.
		return ErrNoHosts
	} else if err = rec.Addr.UnmarshalText(field); err != nil {
		// Don't wrap the error, because it's informative enough as is.
		return err
	}

	// Convert to string prematurely since it seems to be more performant than
	// copying each subslice to a string.  See [BenchmarkRecord_UnmarshalText].
	hosts := string(data)

	n := 0
	for f, t := cutStringField(hosts); f != ""; f, t = cutStringField(t) {
		err = netutil.ValidateDomainName(f)
		if err != nil {
			err = fmt.Errorf("name at index %d: %w", n, err)

			break
		}

		n++
	}

	rec.Names = make([]string, n)
	for i := range rec.Names {
		rec.Names[i], hosts = cutStringField(hosts)
	}

	return err
}

// cutStringField cuts the first substring of data separated by spaces from the
// following part and returns it, trimming leading spaces from tail.  If there
// are no spaces in data, then tail is empty.
func cutStringField(data string) (field, tail string) {
	if endIdx := strings.IndexAny(data, spaces); endIdx < 0 {
		return data, ""
	} else {
		return data[:endIdx], strings.TrimLeft(data[endIdx:], spaces)
	}
}

// cutField cuts the first subslice of data separated by spaces from the
// following part and returns it, trimming leading spaces from tail.  If there
// are no spaces in data, then tail is nil.
func cutField(data []byte) (field, tail []byte) {
	if endIdx := bytes.IndexAny(data, spaces); endIdx < 0 {
		return data, nil
	} else {
		return data[:endIdx], bytes.TrimLeft(data[endIdx:], spaces)
	}
}

// type check
var _ encoding.TextMarshaler = Record{}

// MarshalText implements the [encoding.TextMarshaler] interface for Record.
func (rec Record) MarshalText() (data []byte, err error) {
	namesLen := 0
	for _, name := range rec.Names {
		namesLen += 1 + len(name)
	}

	data, _ = rec.Addr.MarshalText()
	data = slices.Grow(data, namesLen)

	for _, name := range rec.Names {
		data = append(data, ' ')
		data = append(data, name...)
	}

	return slices.Clip(data), nil
}
