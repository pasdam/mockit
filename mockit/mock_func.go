package mockit

import (
	"errors"
	"reflect"
	"testing"

	"bou.ke/monkey"
)

type mockFunc struct {
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

	index, err := m.findCall(newCall.in)
	if err == nil {
		newCall.count = m.calls[index].count
		m.calls[index] = newCall
	} else {
		m.calls = append(m.calls, newCall)
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
	index, err := m.findCall(inValues)
	if err != nil || m.calls[index].count < 1 {
		t.Error("Excepted call didn't happen")
	}
}

func (m *mockFunc) findCall(in []reflect.Value) (int, error) {
	for i := 0; i < len(m.calls); i++ {
		if len(m.calls[i].in) != len(in) {
			break
		}
		if valuesArrayMatch(m.calls[i].in, in) {
			return i, nil
		}
	}

	return -1, errors.New("Unable to find a call with the specified input parameters")
}

func (m *mockFunc) makeCall(in []reflect.Value) []reflect.Value {
	call := m.recordCall(in)
	if call.out != nil {
		return call.out
	}

	return m.defaultOut
}

func (m *mockFunc) recordCall(in []reflect.Value) *call {
	index, err := m.findCall(in)
	if err != nil {
		m.calls = append(m.calls, &call{
			in: in,
		})
		index = len(m.calls) - 1
	}
	m.calls[index].count = m.calls[index].count + 1
	return m.calls[index]
}
