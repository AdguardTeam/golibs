// Logging functions

package log

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
)

const (
	ERROR = iota
	INFO  = iota
	DEBUG = iota
	TRACE = iota
)

// Logging level
var Level uint = INFO

// Write to error log and exit application
func Fatal(v ...interface{}) {
	log.Fatal(v...)
}

// Write to error log and exit application
func Fatalf(format string, v ...interface{}) {
	log.Fatalf(format, v...)
}

// Write to error log
func Error(fmt string, args ...interface{}) {
	log.Printf("[error] "+fmt+"\n", args...)
}

// Write to info log
func Print(v ...interface{}) {
	if Level >= INFO {
		log.Print(v...)
	}
}

// Write to info log
func Printf(format string, v ...interface{}) {
	if Level >= INFO {
		Info(format, v...)
	}
}

// Write to info log
func Println(v ...interface{}) {
	if Level >= INFO {
		log.Println(v...)
	}
}

// Write to info log
func Info(fmt string, args ...interface{}) {
	if Level >= INFO {
		log.Printf("[info] "+fmt+"\n", args...)
	}
}

// Write to debug log
func Debug(fmt string, args ...interface{}) {
	if Level >= DEBUG {
		log.Printf("[debug] "+fmt+"\n", args...)
	}
}

// Write to debug log with the calling function's name
func DebugFunc(format string, args ...interface{}) {
	if Level >= DEBUG {
		pc := make([]uintptr, 10)
		runtime.Callers(2, pc)
		f := runtime.FuncForPC(pc[0])

		var buf strings.Builder
		buf.WriteString(fmt.Sprintf("[debug] %s(): ", f.Name()))
		buf.WriteString(fmt.Sprintf(format+"\n", args...))
		log.Print(buf.String())
	}
}

// Write to trace log with the calling function's name
func Tracef(format string, v ...interface{}) {
	if Level >= TRACE {
		pc := make([]uintptr, 10)
		runtime.Callers(2, pc)
		f := runtime.FuncForPC(pc[0])

		var buf strings.Builder
		buf.WriteString(fmt.Sprintf("%s(): ", path.Base(f.Name())))
		text := fmt.Sprintf(format, v...)
		buf.WriteString(text)
		if len(text) == 0 || text[len(text)-1] != '\n' {
			buf.WriteRune('\n')
		}
		fmt.Fprint(os.Stderr, buf.String())
	}
}
