package mockit

import (
	"reflect"
	"strings"
	"testing"

	"github.com/pasdam/mockit/internal/format"
)

type instanceMock struct {
	defaultOut  []reflect.Value
	calls       [][]reflect.Value
	enabled     bool
	mockedCalls *callsIndex
	t           *testing.T
	target      *reflect.Value
}

func (m *instanceMock) Disable() {
	m.enabled = false
}

func (m *instanceMock) Enable() {
	m.enabled = true
}

func (m *instanceMock) Verify(in ...interface{}) {
	inValues := interfacesArrayToValuesArray(in, m.target.Type().In)
	_, err := findCall(m.calls, inValues, func(fromCalls, in []reflect.Value) bool {
		return callsMatch(in, fromCalls, true)
	})
	if err != nil {
		builder := strings.Builder{}
		builder.WriteString("Expected call: ")
		builder.WriteString(format.PrintCall(m.target, inValues))
		if len(m.calls) > 0 {
			builder.WriteString("; but it recorded the following instead: ")
			builder.WriteString(format.PrintCall(m.target, m.calls[0]))
			for i := 1; i < len(m.calls); i++ {
				builder.WriteString(", ")
				builder.WriteString(format.PrintCall(m.target, m.calls[i]))
			}

		} else {
			builder.WriteString("; but no call was recorded")
		}
		m.t.Error(builder.String())
	}
}

func (m *instanceMock) With(values ...interface{}) Stub {
	typeOf := m.target.Type()
	builder := &stubBuilder{
		args: convertToValuesAndVerifies(m.t, values, typeOf.NumIn(), typeOf.In),
		mock: m,
	}

	return builder
}

func (m *instanceMock) RecordCall(in []reflect.Value) {
	m.calls = append(m.calls, in)
}
