package mockit

import (
	"reflect"
	"testing"

	"bou.ke/monkey"
)

type mockFunc struct {
	funcMockData

	currentMock *funcCall
	defaultOut  []reflect.Value
	guard       *monkey.PatchGuard
}

// MockFunc creates a new Mock to mock a function
func MockFunc(t *testing.T, targetFn interface{}) Mock {
	target := reflect.ValueOf(targetFn)
	if target.Kind() != reflect.Func {
		// TODO: add unknown type to message
		t.Errorf("The target type is not a function, unable to mock it")
		return nil
	}

	mock := &mockFunc{
		defaultOut: defaultFuncOutput(target.Type()),
		funcMockData: funcMockData{
			target: target,
			t:      t,
		},
	}

	replacement := reflect.MakeFunc(reflect.TypeOf(targetFn), mock.makeCall)

	mock.guard = monkey.Patch(targetFn, replacement.Interface())

	t.Cleanup(mock.Disable)

	return mock
}

func (f *mockFunc) CallRealMethod() {
	f.completeMock(nil)
}

func (f *mockFunc) Disable() {
	f.guard.Unpatch()
}

func (f *mockFunc) Enable() {
	f.guard.Restore()
}

func (f *mockFunc) Return(values ...interface{}) {
	typeOf := f.target.Type()
	f.completeMock(f.convertToValuesAndVerifies(values, typeOf.NumOut(), typeOf.Out))
}

func (f *mockFunc) ReturnDefaults() {
	f.completeMock(f.defaultOut)
}

func (f *mockFunc) Verify(in ...interface{}) {
	inValues := interfacesArrayToValuesArray(in, f.target.Type().In)
	_, err := findCall(f.calls, inValues, func(fromCalls, in []reflect.Value) bool {
		return callsMatch(in, fromCalls, true)
	})
	if err != nil {
		f.t.Error("Excepted call didn't happen")
	}
}

func (f *mockFunc) With(values ...interface{}) Stub {
	typeOf := f.target.Type()
	f.currentMock = &funcCall{
		in: f.convertToValuesAndVerifies(values, typeOf.NumIn(), typeOf.In),
	}
	return f
}

func (f *mockFunc) completeMock(out []reflect.Value) {
	f.currentMock.out = out
	f.mocks = append(f.mocks, f.currentMock)
	f.currentMock = nil
}

func (f *mockFunc) convertToValuesAndVerifies(values []interface{}, expectedValuesCount int, expectedValueProvider func(int) reflect.Type) []reflect.Value {
	result := interfacesArrayToValuesArray(values, expectedValueProvider)

	err := verifyValues(expectedValuesCount, expectedValueProvider, result)
	if err != nil {
		f.t.Errorf("Invalid arguments. %s", err.Error())
		return nil
	}

	return result
}

func (f *mockFunc) makeCall(in []reflect.Value) []reflect.Value {
	return makeCall(&f.funcMockData, in, f.defaultOut, f.guard)
}
