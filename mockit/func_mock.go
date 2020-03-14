package mockit

import (
	"reflect"
	"testing"

	"bou.ke/monkey"
)

type funcMock struct {
	mocks       []*call
	calls       []*call
	defaultOut  []reflect.Value
	target      reflect.Value
	guard       *monkey.PatchGuard
	currentMock *call
	t           *testing.T
}

// NewFuncMock creates a new Mock to mock a function
func NewFuncMock(t *testing.T, targetFn interface{}) Mock {
	target := reflect.ValueOf(targetFn)
	if target.Kind() != reflect.Func {
		// TODO: add unknown type to message
		t.Errorf("The target type is not a function, unable to mock it")
		return nil
	}

	mock := &funcMock{
		defaultOut: defaultFuncOutput(target.Type()),
		target:     target,
	}

	replacement := reflect.MakeFunc(reflect.TypeOf(targetFn), mock.makeCall)

	mock.guard = monkey.Patch(targetFn, replacement.Interface())

	return mock
}

func (f *funcMock) CallRealMethod() {
	f.completeMock(nil)
}

func (f *funcMock) Disable() {
	f.guard.Unpatch()
}

func (f *funcMock) Enable() {
	f.guard.Restore()
}

func (f *funcMock) Return(values []interface{}) {
	typeOf := f.target.Type()
	f.completeMock(f.convertToValuesAndVerifies(values, typeOf.NumOut(), typeOf.Out))
}

func (f *funcMock) ReturnDefaults() {
	f.completeMock(f.defaultOut)
}

func (f *funcMock) Verify(t *testing.T, in []interface{}) {
	inValues := interfacesArrayToValuesArray(in, f.target.Type().In)
	_, err := findCall(f.calls, inValues, func(fromCalls, in []reflect.Value) bool {
		return callsMatch(in, fromCalls, true)
	})
	if err != nil {
		t.Error("Excepted call didn't happen")
	}
}

func (f *funcMock) With(values []interface{}) Stub {
	typeOf := f.target.Type()
	f.currentMock = &call{
		in: f.convertToValuesAndVerifies(values, typeOf.NumIn(), typeOf.In),
	}
	return f
}

func (f *funcMock) completeMock(out []reflect.Value) {
	f.currentMock.out = out
	f.mocks = append(f.mocks, f.currentMock)
	f.currentMock = nil
}

func (f *funcMock) convertToValuesAndVerifies(values []interface{}, expectedValuesCount int, expectedValueProvider func(int) reflect.Type) []reflect.Value {
	result := interfacesArrayToValuesArray(values, expectedValueProvider)

	err := verifyValues(expectedValuesCount, expectedValueProvider, result)
	if err != nil {
		f.t.Errorf("Invalid arguments. %s", err.Error())
		return nil
	}

	return result
}

func (f *funcMock) makeCall(in []reflect.Value) []reflect.Value {
	// record call
	f.calls = append(f.calls, &call{
		in: in,
	})

	index, err := findCall(f.mocks, in, func(fromCalls, in []reflect.Value) bool {
		return callsMatch(fromCalls, in, true)
	})
	if err == nil { // mock exists
		if f.mocks[index].out == nil {
			// disable mock
			f.guard.Unpatch()
			defer f.guard.Restore()

			// call real method
			return f.target.Call(in)
		}
		// return expected values
		return f.mocks[index].out
	}

	return f.defaultOut
}
