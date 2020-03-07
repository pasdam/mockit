package mockit

import (
	"reflect"
)

func interfacesArrayToValuesArray(args []interface{}, defaultValueProvider func(int) reflect.Type) []reflect.Value {
	result := make([]reflect.Value, 0, len(args))

	for i := 0; i < len(args); i++ {
		if args[i] == nil {
			result = append(result, reflect.Zero(defaultValueProvider(i)))
		} else {
			result = append(result, reflect.ValueOf(args[i]))
		}
	}

	return result
}
