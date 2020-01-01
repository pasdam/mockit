package mockit

import (
	"reflect"
)

func valuesArrayMatch(expected []reflect.Value, actual []reflect.Value) bool {
	if len(expected) != len(actual) {
		return false
	}
	for i := 0; i < len(expected); i++ {
		if !reflect.DeepEqual(expected[i].Interface(), actual[i].Interface()) {
			return false
		}
	}
	return true
}
