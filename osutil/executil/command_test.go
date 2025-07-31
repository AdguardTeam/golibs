package executil_test

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/AdguardTeam/golibs/osutil/executil"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSystemCommand_Shutdown(t *testing.T) {
	cons := executil.SystemCommandConstructor{}

	const sleepTime = 5 * time.Second
	ctx := testutil.ContextWithTimeout(t, sleepTime*2)

	var sleepCmd string
	if runtime.GOOS == "windows" {
		sleepCmd = fmt.Sprintf("Start-Sleep %v", sleepTime.Seconds())
	} else {
		sleepCmd = fmt.Sprintf("sleep %v", sleepTime.Seconds())
	}

	c, err := cons.New(ctx, &executil.CommandConfig{
		Path: shell,
		Args: []string{"-c", sleepCmd},
	})
	require.NoError(t, err)

	err = c.Start(testutil.ContextWithTimeout(t, sleepTime/10))
	require.NoError(t, err)

	go func() {
		sdErr := c.Cancel(testutil.ContextWithTimeout(t, sleepTime/10))
		assert.NoError(t, sdErr)
	}()

	err = c.Wait(testutil.ContextWithTimeout(t, sleepTime/10))
	require.Error(t, err)

	code, ok := executil.ExitCodeFromError(err)
	require.True(t, ok)

	assert.NotEqual(t, 0, code)
}
