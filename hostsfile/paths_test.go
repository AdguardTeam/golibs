package hostsfile_test

import (
	"testing"

	"github.com/AdguardTeam/golibs/hostsfile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultHostsPaths(t *testing.T) {
	t.Parallel()

	paths, err := hostsfile.DefaultHostsPaths()
	require.NoError(t, err)

	assert.NotEmpty(t, paths)
}
