package mockit

import (
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/pasdam/mockit/internal/utils"
)

type mockMethodGuard struct {
	defaultOut []reflect.Value
	guard      *monkey.PatchGuard
	methodName string
	mocks      map[interface{}]*mockMethod
	t          *testing.T
}

func newMockMethodGuard(t *testing.T, fullyQualifiedName string, method reflect.Value, instance interface{}) *mockMethodGuard {
	methodName := utils.MethodName(fullyQualifiedName)

	instanceType := reflect.TypeOf(instance)

	methodType, found := instanceType.MethodByName(methodName)
	if !found {
		t.Errorf("The specified instance does not have a method called %s", methodName)
		return nil
	}

	mock := &mockMethodGuard{
		defaultOut: defaultFuncOutput(methodType.Type),
		methodName: methodName,
		mocks:      make(map[interface{}]*mockMethod),
		t:          t,
	}

	replacement := reflect.MakeFunc(methodType.Func.Type(), mock.makeCall)
	mock.guard = monkey.PatchInstanceMethod(instanceType, methodName, replacement.Interface())

	return mock
}

func (m *mockMethodGuard) methodMock(instance interface{}) *mockMethod {
	mock, found := m.mocks[instance]
	if !found {
		mock = newMethodMock(m.t, instance)
		mock.defaultOut = m.defaultOut
		m.mocks[instance] = mock
	}
	return mock
}

func (m *mockMethodGuard) makeCall(in []reflect.Value) []reflect.Value {
	instanceValue := in[0]
	instance := instanceValue.Interface()
	mock, found := m.mocks[instance]
	if !found || !mock.enabled {
		m.guard.Unpatch()
		defer m.guard.Restore()
		return instanceValue.MethodByName(m.methodName).Call(in[1:])
	}

	return makeCall(&mock.funcMockData, in[1:], m.defaultOut, m.guard)
}
