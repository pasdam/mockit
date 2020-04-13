package mockit

import (
	"reflect"
)

func verifyCall(f *funcMockData, in ...interface{}) {
	inValues := interfacesArrayToValuesArray(in, f.target.Type().In)
	_, err := findCall(f.calls, inValues, func(fromCalls, in []reflect.Value) bool {
		return callsMatch(in, fromCalls, true)
	})
	if err != nil {
		f.t.Error("Expected call didn't happen")
	}
}
