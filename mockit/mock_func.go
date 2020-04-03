package mockit

import (
	"reflect"
	"testing"

	"bou.ke/monkey"
)

type mockFunc struct {
	funcMockData

	guard *monkey.PatchGuard
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
		funcMockData: funcMockData{
			defaultOut: defaultFuncOutput(target.Type()),
			target:     target,
			t:          t,
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
	configureMockReturn(&f.funcMockData, values...)
}

func (f *mockFunc) ReturnDefaults() {
	f.completeMock(f.defaultOut)
}

func (f *mockFunc) Verify(in ...interface{}) {
	verifyCall(&f.funcMockData, in...)
}

func (f *mockFunc) With(values ...interface{}) Stub {
	configureMockWith(&f.funcMockData, values...)
	return f
}

func (f *mockFunc) makeCall(in []reflect.Value) []reflect.Value {
	return makeCall(&f.funcMockData, in, f.defaultOut, f.guard)
}
