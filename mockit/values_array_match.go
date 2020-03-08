package mockit

import (
	"reflect"

	"github.com/pasdam/mockit/matchers/argument"
)

func valuesArrayMatch(expected []reflect.Value, actual []reflect.Value) bool {
	if len(expected) != len(actual) {
		return false
	}

	var matcher argument.Matcher
	matcherType := reflect.TypeOf(matcher)

	for i := 0; i < len(expected); i++ {
		if expected[i].Type().AssignableTo(matcherType) {
			return expected[i].Call([]reflect.Value{actual[i]})[0].Bool()
		}
		if !reflect.DeepEqual(expected[i].Interface(), actual[i].Interface()) {
			return false
		}
	}
	return true
}
