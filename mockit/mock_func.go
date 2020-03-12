package mockit

import (
	"errors"
	"reflect"
	"testing"

	"bou.ke/monkey"
)

type mockFunc struct {
	mocks      []*call
	calls      []*call
	defaultOut []reflect.Value
	typeOf     reflect.Type
	guard      *monkey.PatchGuard
}

// NewMockFunc creates a new function mock instance
func NewMockFunc(t *testing.T, target interface{}) Mock {
	typeOf := reflect.TypeOf(target)
	if typeOf.Kind() != reflect.Func {
		// TODO: add unknown type to message
		t.Errorf("The target type is not a function, unable to mock it")
		return nil
	}

	mock := &mockFunc{
		defaultOut: defaultFuncOutput(typeOf),
		typeOf:     typeOf,
	}

	replacement := reflect.MakeFunc(reflect.TypeOf(target), mock.makeCall)

	mock.guard = monkey.Patch(target, replacement.Interface())

	return mock
}

func (m *mockFunc) Mock(t *testing.T, in []interface{}, out []interface{}) {
	newCall := &call{
		in:  interfacesArrayToValuesArray(in, m.typeOf.In),
		out: interfacesArrayToValuesArray(out, m.typeOf.Out),
	}

	err := verifyValues(m.typeOf.NumIn(), m.typeOf.In, newCall.in)
	if err != nil {
		t.Errorf("Invalid input. %s", err.Error())
		return
	}
	err = verifyValues(m.typeOf.NumOut(), m.typeOf.Out, newCall.out)
	if err != nil {
		t.Errorf("Invalid output. %s", err.Error())
		return
	}

	index, err := findCall(m.mocks, newCall.in, func(fromCalls, in []reflect.Value) bool {
		return callsMatch(fromCalls, in, false)
	})
	if err == nil {
		m.mocks[index] = newCall
	} else {
		m.mocks = append(m.mocks, newCall)
	}
}

func (m *mockFunc) Restore() {
	m.guard.Restore()
}

func (m *mockFunc) UnMock() {
	m.guard.Unpatch()
}

func (m *mockFunc) Verify(t *testing.T, in []interface{}) {
	inValues := interfacesArrayToValuesArray(in, m.typeOf.In)
	_, err := findCall(m.calls, inValues, func(fromCalls, in []reflect.Value) bool {
		return callsMatch(in, fromCalls, true)
	})
	if err != nil {
		t.Error("Excepted call didn't happen")
	}
}

// Mock: calls=m.mocks, in=mock.in
// Verify: calls=m.calls, in=expectedIn

func findCall(calls []*call, in []reflect.Value, matcher func(fromCalls, in []reflect.Value) bool) (int, error) {
	for i := 0; i < len(calls); i++ {
		if matcher(calls[i].in, in) {
			return i, nil
		}
	}

	return -1, errors.New("Unable to find a call with the specified input parameters")
}

func (m *mockFunc) makeCall(in []reflect.Value) []reflect.Value {
	m.recordCall(in)

	index, err := findCall(m.mocks, in, func(fromCalls, in []reflect.Value) bool {
		return callsMatch(fromCalls, in, true)
	})
	if err == nil {
		return m.mocks[index].out
	}

	return m.defaultOut
}

func (m *mockFunc) recordCall(in []reflect.Value) {
	m.calls = append(m.calls, &call{
		in: in,
	})
}
