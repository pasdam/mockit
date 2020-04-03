package mockit

import (
	"reflect"
	"testing"
)

type funcMockData struct {
	calls       []*funcCall
	currentMock *funcCall
	instance    interface{}
	mocks       []*funcCall
	t           *testing.T
	target      reflect.Value
	defaultOut  []reflect.Value
}

func (f *funcMockData) completeMock(out []reflect.Value) {
	f.currentMock.out = out
	f.mocks = append(f.mocks, f.currentMock)
	f.currentMock = nil
}

func (f *funcMockData) convertToValuesAndVerifies(values []interface{}, expectedValuesCount int, expectedValueProvider func(int) reflect.Type) []reflect.Value {
	result := interfacesArrayToValuesArray(values, expectedValueProvider)

	err := verifyValues(expectedValuesCount, expectedValueProvider, result)
	if err != nil {
		f.t.Errorf("Invalid arguments. %s", err.Error())
		return nil
	}

	return result
}
