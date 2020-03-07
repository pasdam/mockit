package mockit

import (
	"reflect"
)

func valuesToInterfaceArray(args []reflect.Value) []interface{} {
	emptyVal := reflect.Value{}
	result := make([]interface{}, 0, len(args))

	for _, arg := range args {
		if arg == emptyVal {
			result = append(result, nil)

		} else {
			result = append(result, arg.Interface())
		}
	}

	return result
}
