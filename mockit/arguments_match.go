package mockit

import (
	"reflect"

	"github.com/pasdam/mockit/matchers/argument"
)

var matcherType = reflect.TypeOf(argument.Any)

func argumentsMatch(expected reflect.Value, actual reflect.Value, enableMatcher bool) bool {
	equal := reflect.DeepEqual(expected.Interface(), actual.Interface())
	if equal {
		return true
	}

	if enableMatcher && expected.Type().AssignableTo(matcherType) {
		return expected.Call([]reflect.Value{actual})[0].Bool()
	}

	return false
}
