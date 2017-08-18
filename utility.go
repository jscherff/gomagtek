package gomagtek

import "path/filepath"
import "runtime"
import "fmt"

// getFunctionInfo returns function filename, line number, and function name
// for error reporting.
func getFunctionInfo() string {

	pc, file, line, success := runtime.Caller(1)
	function := runtime.FuncForPC(pc)

	if !success {
		return "Unknown goroutine"
	}

	return fmt.Sprintf("%s:%d: %s()", filepath.Base(file), line, function.Name())
}
