package hostsfile_test

import (
	"time"

	"github.com/AdguardTeam/golibs/logutil/slogutil"
)

// testTimeout is a common timeout for tests.
const testTimeout = 1 * time.Second

// testLogger is a common logger for tests.
var testLogger = slogutil.NewDiscardLogger()
