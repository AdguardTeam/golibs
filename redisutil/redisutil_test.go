package redisutil_test

import (
	"time"

	"github.com/AdguardTeam/golibs/logutil/slogutil"
)

// testTimeout is the common timeout for tests.
const testTimeout = 1 * time.Second

// Key and value constants.
const (
	testKey   = "key"
	testValue = "value"
)

// testLogger is the common logger for tests.
var testLogger = slogutil.NewDiscardLogger()
