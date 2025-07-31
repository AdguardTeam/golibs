package executil_test

import (
	"context"
	"fmt"
	"runtime"

	"github.com/AdguardTeam/golibs/osutil/executil"
	"github.com/c2h5oh/datasize"
)

func ExampleRunWithPeek() {
	cons := executil.SystemCommandConstructor{}

	ctx := context.Background()

	okCmd := `exit 0`
	errCmd := `printf 'a long error!\n' 1>&2 && exit 1`
	if runtime.GOOS == "windows" {
		errCmd = `$host.ui.WriteErrorLine('a long error!') ; exit 1`
	}

	const limit = 8 * datasize.B
	err := executil.RunWithPeek(ctx, cons, limit, shell, "-c", okCmd)
	fmt.Println(err)

	err = executil.RunWithPeek(ctx, cons, limit, shell, "-c", errCmd)
	fmt.Println(err)

	// Output:
	// <nil>
	// running: exit status 1; stderr peek: "a long e"; stdout peek: ""
}
