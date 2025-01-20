package service_test

import (
	"context"
	"log/slog"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/logutil/slogutil"
	"github.com/AdguardTeam/golibs/service"
)

func ExampleSlogErrorHandler() {
	logger := slogutil.New(&slogutil.Config{
		Format: slogutil.FormatJSON,
	})

	h := service.NewSlogErrorHandler(logger, slog.LevelError, "test message")

	h.Handle(context.Background(), errors.Error("test error"))

	// Output:
	// {"level":"ERROR","msg":"test message","err":"test error"}
}
