package mockit

import (
	"fmt"
	"reflect"
)

func verifyValues(expectedCount int, expectedValueProvider func(int) reflect.Type, actualValues []reflect.Value) error {
	if expectedCount != len(actualValues) {
		return fmt.Errorf("Expected values count (%d) is different than the actual size (%d)", expectedCount, len(actualValues))
	}

	for i := 0; i < expectedCount; i++ {
		expected := expectedValueProvider(i)
		actual := actualValues[i].Type()
		if expected != actual && !actual.AssignableTo(expected) {
			return fmt.Errorf("Type at index %d is different than expected (%v): actual type %v", i, expectedValueProvider(i), actualValues[i].Type())
		}
	}

	return nil
}
