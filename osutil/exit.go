package osutil

import "os"

// ExitFunc represents a function that is used to exit the application.  Note,
// that [os.Exit] has the same signature.
type ExitFunc func(code int)

// type check
var _ ExitFunc = os.Exit
