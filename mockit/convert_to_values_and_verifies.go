package mockit

import (
	"reflect"
	"testing"
)

func convertToValuesAndVerifies(t *testing.T, values []interface{}, expectedValuesCount int, expectedValueProvider func(int) reflect.Type) []reflect.Value {
	result := interfacesArrayToValuesArray(values, expectedValueProvider)

	err := verifyValues(expectedValuesCount, expectedValueProvider, result)
	if err != nil {
		t.Errorf("Invalid arguments. %s", err.Error())
		return nil
	}

	return result
}
