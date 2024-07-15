package httputil_test

import (
	"time"

	"github.com/AdguardTeam/golibs/netutil/httputil"
)

// testTimeout is a common timeout for tests.
const testTimeout = 1 * time.Second

// Common constants for tests.
const (
	testPath = "/health-check"
	testBody = string(httputil.HealthCheckHandler)
)
