package ioutil_test

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/AdguardTeam/golibs/ioutil"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLimitedReader_Read(t *testing.T) {
	testCases := []struct {
		err   error
		name  string
		in    string
		limit uint64
		want  int
	}{{
		err:   nil,
		name:  "perfectly_match",
		in:    "abc",
		limit: 3,
		want:  3,
	}, {
		err:   io.EOF,
		name:  "eof",
		in:    "",
		limit: 3,
		want:  0,
	}, {
		err: &ioutil.LimitError{
			Limit: 0,
		},
		name:  "limit_reached",
		in:    "abc",
		limit: 0,
		want:  0,
	}, {
		err:   nil,
		name:  "truncated",
		in:    "abc",
		limit: 2,
		want:  2,
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			readCloser := io.NopCloser(strings.NewReader(tc.in))
			buf := make([]byte, tc.limit+1)

			limited := ioutil.LimitReader(readCloser, tc.limit)
			n, err := limited.Read(buf)
			require.Equal(t, tc.err, err)

			assert.Equal(t, tc.want, n)
		})
	}
}

func TestLimitError_Error(t *testing.T) {
	const limit = 42
	err := &ioutil.LimitError{
		Limit: limit,
	}

	want := fmt.Sprintf("cannot read more than %d bytes", limit)
	testutil.AssertErrorMsg(t, want, err)
}
