package mockit

import (
	"reflect"
)

func interfacesArrayToValuesArray(args []interface{}) []reflect.Value {
	result := make([]reflect.Value, 0, len(args))

	for _, arg := range args {
		result = append(result, reflect.ValueOf(arg))
	}

	return result
}
