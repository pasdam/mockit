package mockit

import (
	"reflect"
	"testing"

	"github.com/pasdam/mockit/internal/utils"
)

var methodMocksMap = make(map[string]*mockMethodGuard)

type mockMethod struct {
	funcMockData

	enabled bool
}

// MockMethod creates a new Mock to mock an instance method
func MockMethod(t *testing.T, instance interface{}, method interface{}) Mock {
	if instance == nil {
		t.Errorf("Instance can't be nil")
		return nil
	}
	if method == nil {
		t.Errorf("Method can't be nil")
		return nil
	}

	methodValue := reflect.ValueOf(method)
	if methodValue.Type().Kind() != reflect.Func {
		t.Errorf("The 'method' parameter is not actually a func")
		return nil
	}

	fullyQualifiedName := utils.MethodFullyQualifiedName(methodValue)

	methodMockGuard, ok := methodMocksMap[fullyQualifiedName]
	if !ok {
		methodMockGuard = newMockMethodGuard(t, fullyQualifiedName, methodValue, instance)
		methodMocksMap[fullyQualifiedName] = methodMockGuard
	}

	mock := methodMockGuard.methodMock(instance)
	mock.target = methodValue

	return mock
}

func newMethodMock(t *testing.T, instance interface{}) *mockMethod {
	mock := &mockMethod{
		enabled: true,
		funcMockData: funcMockData{
			instance: instance,
			t:        t,
		},
	}

	t.Cleanup(mock.Disable)

	return mock
}

func (m *mockMethod) CallRealMethod() {
	m.completeMock(nil)
}

func (m *mockMethod) Disable() {
	m.enabled = false
}

func (m *mockMethod) Enable() {
	m.enabled = true
}

func (m *mockMethod) Return(values ...interface{}) {
	configureMockReturn(&m.funcMockData, values...)
}

func (m *mockMethod) ReturnDefaults() {
	m.completeMock(m.defaultOut)
}

func (m *mockMethod) Verify(in ...interface{}) {
	verifyCall(&m.funcMockData, in...)
}

func (m *mockMethod) With(values ...interface{}) Stub {
	configureMockWith(&m.funcMockData, values...)
	return m
}
