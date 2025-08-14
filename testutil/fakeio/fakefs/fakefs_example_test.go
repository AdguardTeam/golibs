package fakefs_test

import (
	"fmt"
	"io"
	"io/fs"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/AdguardTeam/golibs/testutil/fakeio/fakefs"
)

func Example() {
	const testMsg = "test message"

	isClosed := false
	fakeFile := &fakefs.File{
		// Use OnClose to record the closure of the file.
		OnClose: func() (err error) {
			isClosed = true

			return nil
		},

		// Use OnRead to return the fake data.
		OnRead: func(b []byte) (n int, err error) {
			copy(b, testMsg)

			return len(testMsg), io.EOF
		},

		// Use OnStat with a panic to signal that Stat is expected to not be
		// called.
		OnStat: func() (fi fs.FileInfo, err error) { panic(testutil.UnexpectedCall()) },
	}

	fakeFS := &fakefs.FS{
		OnOpen: func(_ string) (f fs.File, err error) {
			return fakeFile, nil
		},
	}

	// The function that is expected to call Read and Close.
	testedFunction := func(fsys fs.FS) (b []byte, err error) {
		f, err := fsys.Open("my_file.txt")
		if err != nil {
			return nil, fmt.Errorf("opening: %w", err)
		}
		defer func() { err = errors.WithDeferred(err, f.Close()) }()

		b = make([]byte, len(testMsg)*2)
		n, err := f.Read(b)
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("reading: %w", err)
		}

		return b[:n], nil
	}

	// A simulation of a successful test.
	gotData, gotErr := testedFunction(fakeFS)
	fmt.Printf("read: %v %q\n", gotErr, gotData)
	fmt.Printf("closed: %t\n", isClosed)

	// Output:
	// read: <nil> "test message"
	// closed: true
}
