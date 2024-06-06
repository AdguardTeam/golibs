package slogutil_test

import (
	"testing"

	"github.com/AdguardTeam/golibs/logutil/slogutil"
	"github.com/stretchr/testify/assert"
)

func TestNewFormat(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		f, err := slogutil.NewFormat(string(slogutil.FormatJSON))
		assert.NotEmpty(t, f)
		assert.Nil(t, err)
	})

	t.Run("bad", func(t *testing.T) {
		const badFmt = "not a real format"

		f, err := slogutil.NewFormat(badFmt)
		assert.Empty(t, f)
		assert.Equal(t, err, &slogutil.BadFormatError{
			Format: badFmt,
		})
	})
}
