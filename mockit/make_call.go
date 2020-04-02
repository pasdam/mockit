package mockit

import (
	"reflect"

	"bou.ke/monkey"
)

func makeCall(mock *funcMockData, in []reflect.Value, defaultOut []reflect.Value, guard *monkey.PatchGuard) []reflect.Value {
	// record call
	mock.calls = append(mock.calls, &funcCall{in: in})

	index, err := findCall(mock.mocks, in, func(fromCalls, in []reflect.Value) bool {
		return callsMatch(fromCalls, in, true)
	})
	if err == nil { // mock exists
		if mock.mocks[index].out == nil {
			// disable mock
			guard.Unpatch()
			defer guard.Restore()

			// call real method
			return mock.target.Call(in)
		}
		// return expected values
		return mock.mocks[index].out
	}

	return defaultOut
}
