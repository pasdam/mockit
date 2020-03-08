package mockit

import (
	"fmt"
	"reflect"

	"github.com/pasdam/mockit/matchers/argument"
)

func verifyValues(expectedCount int, expectedValueProvider func(int) reflect.Type, actualValues []reflect.Value) error {
	if expectedCount != len(actualValues) {
		return fmt.Errorf("Expected values count (%d) is different than the actual size (%d)", expectedCount, len(actualValues))
	}

	var matcher argument.Matcher
	matcherType := reflect.TypeOf(matcher)
	emptyVal := reflect.Value{}

	for i := 0; i < expectedCount; i++ {
		expected := expectedValueProvider(i)
		actualValue := actualValues[i]
		if actualValue == emptyVal {
			if canBeNil(expected.Kind()) {
				continue

			} else {
				return fmt.Errorf("Cannot assign nil at index %d to the type %v", i, expected)
			}

		} else if actualValue.Type().AssignableTo(matcherType) {
			continue
		}

		actual := actualValue.Type()
		if expected != actual && !actual.AssignableTo(expected) {
			return fmt.Errorf("Type at index %d is different than expected (%v): actual type %v", i, expected, actualValues[i].Type())
		}
	}

	return nil
}
