package slogutil_test

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/logutil/slogutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecoverAndLog(t *testing.T) {
	const (
		errTestMsg              = "test error"
		errTest    errors.Error = errTestMsg
	)

	require.True(t, t.Run("no_panic", func(t *testing.T) {
		output := &bytes.Buffer{}
		l := slogutil.New(&slogutil.Config{
			Output: output,
		})

		func() {
			ctx := context.Background()

			defer slogutil.RecoverAndLog(ctx, l)
		}()

		assert.Equal(t, 0, output.Len())
	}))

	require.True(t, t.Run("non_error", func(t *testing.T) {
		output := &bytes.Buffer{}
		l := slogutil.New(&slogutil.Config{
			Output: output,
		})

		func() {
			ctx := context.Background()

			defer slogutil.RecoverAndLog(ctx, l)

			panic(errTestMsg)
		}()

		assertPanic(t, output, fmt.Sprintf("ERROR recovered from panic value=%q", errTestMsg))
	}))

	require.True(t, t.Run("error", func(t *testing.T) {
		output := &bytes.Buffer{}
		l := slogutil.New(&slogutil.Config{
			Output: output,
		})

		func() {
			ctx := context.Background()

			defer slogutil.RecoverAndLog(ctx, l)

			panic(errTest)
		}()

		assertPanic(t, output, fmt.Sprintf("ERROR recovered from panic err=%q", errTest))
	}))
}

// assertPanic is a test helper that checks the panic message.
func assertPanic(tb testing.TB, output *bytes.Buffer, wantFirstLine string) {
	tb.Helper()

	lines := strings.Split(output.String(), "\n")

	// Require that there are at least the first line, the last empty line, and
	// at least one stack trace record, which takes two lines.
	require.Greater(tb, len(lines), 4)

	assert.Equal(tb, wantFirstLine, lines[0])

	// Remove the first line, which has already been inspected, and the last
	// line, since it's likely empty.
	lines = lines[1 : len(lines)-1]

	wantRE := regexp.MustCompilePOSIX(`^ERROR stack i=[0-9]+ line="[^"]+"`)
	for i, line := range lines {
		assert.Regexpf(tb, wantRE, line, "at index %d", i)
	}
}
