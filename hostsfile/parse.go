package hostsfile

import (
	"bufio"
	"context"
	"fmt"
	"io"

	"github.com/AdguardTeam/golibs/errors"
)

// Parse reads src and parses it as a hosts file line by line using buf for
// buffered scanning.  If src is a [NamedReader], the name of the data source
// will be set to the Source field of each [Record].
//
// dst must not be nil, use [DiscardSet] if only the unmarshaling errors needed.
// By default it returns all unmarshaling errors within err, but if dst is also
// a [HandleSet], it will be used to handle invalid records and unmarshaling
// errors wrapped with [LineError], see [Record.UnmarshalText] for returned
// errors.
func Parse(ctx context.Context, dst Set, src io.Reader, buf []byte) (err error) {
	var srcName string
	nr, ok := src.(NamedReader)
	if ok {
		srcName = nr.Name()
	}

	var errs []error
	// By default, collect all errors.
	handleInvalid := func(_ context.Context, _ string, _ []byte, err error) { errs = append(errs, err) }

	if handleSet, isHandleSet := dst.(HandleSet); isHandleSet {
		handleInvalid = handleSet.HandleInvalid
	}

	s := bufio.NewScanner(src)
	s.Buffer(buf, bufio.MaxScanTokenSize)

	// TODO(f.setrakov): Implement a stop on context cancel.
	for lineNum := 1; s.Scan(); lineNum++ {
		data := s.Bytes()
		rec := &Record{Source: srcName}

		err = rec.UnmarshalText(data)
		if err != nil {
			handleInvalid(ctx, srcName, data, &LineError{Line: lineNum, err: err})
		} else {
			dst.Add(ctx, rec)
		}
	}
	if err = s.Err(); err != nil {
		return fmt.Errorf("scanning: %w", err)
	}

	return errors.Annotate(errors.Join(errs...), "parsing: %w")
}
