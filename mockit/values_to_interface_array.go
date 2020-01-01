package mockit

import (
	"reflect"
)

func valuesToInterfaceArray(args []reflect.Value) []interface{} {
	result := make([]interface{}, 0, len(args))

	for _, arg := range args {
		result = append(result, arg.Interface())
	}

	return result
}
