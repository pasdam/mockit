package mockit

import "reflect"

func callsMatch(expectedArgs []reflect.Value, actualArgs []reflect.Value, enableMatchers bool) bool {
	if len(expectedArgs) != len(actualArgs) {
		return false
	}

	for i := 0; i < len(expectedArgs); i++ {
		if !argumentsMatch(expectedArgs[i], actualArgs[i], enableMatchers) {
			return false
		}
	}

	return true
}
