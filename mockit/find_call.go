package mockit

import (
	"errors"
	"reflect"
)

func findCall(calls []*funcCall, in []reflect.Value, matcher func(fromCalls, in []reflect.Value) bool) (int, error) {
	for i := 0; i < len(calls); i++ {
		if matcher(calls[i].in, in) {
			return i, nil
		}
	}

	return -1, errors.New("Unable to find a call with the specified input parameters")
}
