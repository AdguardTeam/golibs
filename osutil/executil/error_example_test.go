package executil_test

import (
	"context"
	"fmt"

	"github.com/AdguardTeam/golibs/osutil/executil"
)

func ExampleExitCodeFromError() {
	cons := executil.SystemCommandConstructor{}

	ctx := context.Background()

	err := executil.Run(ctx, cons, &executil.CommandConfig{
		Path: shell,
		Args: []string{"-c", "exit 123"},
	})
	fmt.Println(err)
	fmt.Println(executil.ExitCodeFromError(err))

	// Output:
	// running: exit status 123
	// 123 true
}
