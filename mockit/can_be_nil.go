package mockit

import (
	"reflect"
)

func canBeNil(kind reflect.Kind) bool {
	switch kind {
	case reflect.Array:
		return true
	case reflect.Chan:
		return true
	case reflect.Func:
		return true
	case reflect.Interface:
		return true
	case reflect.Map:
		return true
	case reflect.Ptr:
		return true
	case reflect.Slice:
		return true
	case reflect.UnsafePointer:
		return true
	}
	return false
}
