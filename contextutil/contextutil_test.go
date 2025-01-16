package contextutil_test

import (
	"context"
	"testing"
	"time"

	"github.com/AdguardTeam/golibs/contextutil"
	"github.com/stretchr/testify/assert"
)

func TestTimeoutConstructor(t *testing.T) {
	const timeout = 1 * time.Minute

	c := contextutil.NewTimeoutConstructor(timeout)
	ctx, cancel := c.New(context.Background())
	defer cancel()

	dl, ok := ctx.Deadline()
	assert.True(t, ok)

	d := time.Until(dl)
	assert.InDelta(t, timeout, d, float64(1*time.Second))
}
