package hostsfile

import (
	"bytes"
	"fmt"

	"github.com/AdguardTeam/golibs/netutil"
)

// UmarshalTextEachSublice is only exists in purposes of benchmarking.
// Currently it demonstrates that premature conversion to string shows better
// performance than conversion to string of every subslice.
//
// See [BenchmarkRecord_UnmarshalText].
func (rec *Record) UnmarshalTextEachSublice(data []byte) (err error) {
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

	n := 0
	for f, t := cutField(data); len(f) > 0; f, t = cutField(t) {
		n++
	}

	rec.Names = make([]string, 0, n)
	for len(data) > 0 {
		field, data = cutField(data)
		name := string(field)
		err = netutil.ValidateDomainName(name)
		if err != nil {
			return fmt.Errorf("name at index %d: %w", len(rec.Names), err)
		}

		rec.Names = append(rec.Names, name)
	}

	return err
}
