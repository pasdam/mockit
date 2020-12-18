package mockit

import (
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/pasdam/mockit/internal/utils"
)

var manager = &mockManager{
	mockedTypes: make(map[string]*mockGuard),
}

type mockManager struct {
	mockedTypes map[string]*mockGuard
}

type patcherProvider func(guard *mockGuard) func(instance interface{}) (*monkey.PatchGuard, callMetadataProvider)

func (m *mockManager) mock(t *testing.T, any bool, instance interface{}, targetFn interface{}, provider patcherProvider) Mock {
	if targetFn == nil {
		t.Error("Method can't be nil")
		return nil
	}

	target := reflect.ValueOf(targetFn)
	if target.Kind() != reflect.Func {
		t.Errorf("The target type (%v) is not a function, unable to mock it", target.Kind())
		return nil
	}

	fullyQualifiedName := utils.MethodFullyQualifiedName(target)

	guard, found := m.mockedTypes[fullyQualifiedName]
	if !found {
		guard = &mockGuard{
			defaultOut:         defaultFuncOutput(target.Type()),
			fullyQualifiedName: fullyQualifiedName,
			mockedInstances:    make(map[interface{}]*instanceMock),
			targetFunc:         target,
		}
		m.mockedTypes[fullyQualifiedName] = guard

		guard.guard, guard.provider = provider(guard)(instance)

		t.Cleanup(func() {
			guard.guard.Unpatch()
			delete(m.mockedTypes, fullyQualifiedName)
		})
	}

	var key interface{}
	if !any {
		key = instance
	}

	mock := guard.mockedInstances[key]
	if mock == nil {
		mock = &instanceMock{
			calls:       nil,
			defaultOut:  guard.defaultOut,
			enabled:     true,
			mockedCalls: &callsIndex{},
			t:           t,
			target:      &target,
		}
		guard.mockedInstances[key] = mock
	}

	return mock
}

func (m *mockManager) MockFunc(t *testing.T, targetFn interface{}) Mock {
	provider := func(guard *mockGuard) func(instance interface{}) (*monkey.PatchGuard, callMetadataProvider) {
		return guard.patchFunc
	}
	return m.mock(t, false, nil, targetFn, provider)
}

func (m *mockManager) MockMethod(t *testing.T, instance interface{}, targetFn interface{}) Mock {
	provider := func(guard *mockGuard) func(instance interface{}) (*monkey.PatchGuard, callMetadataProvider) {
		return guard.patchMethod
	}
	return m.mock(t, false, instance, targetFn, provider)
}

func (m *mockManager) MockMethodForAll(t *testing.T, instance interface{}, targetFn interface{}) Mock {
	provider := func(guard *mockGuard) func(instance interface{}) (*monkey.PatchGuard, callMetadataProvider) {
		return guard.patchMethod
	}
	return m.mock(t, true, instance, targetFn, provider)
}
