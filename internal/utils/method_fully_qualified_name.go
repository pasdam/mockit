package utils

import (
	"reflect"
	"runtime"
)

// MethodFullyQualifiedName returns the fully qualified name of the
// specified method
func MethodFullyQualifiedName(methodValue reflect.Value) string {
	return runtime.FuncForPC(methodValue.Pointer()).Name()
}
