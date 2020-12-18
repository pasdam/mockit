package mockit

import (
	"log"
	"reflect"

	"bou.ke/monkey"
	"github.com/pasdam/mockit/internal/utils"
)

type callMetadataProvider func(in []reflect.Value) (interface{}, *reflect.Value, []reflect.Value)

type mockGuard struct {
	defaultOut         []reflect.Value
	guard              *monkey.PatchGuard
	fullyQualifiedName string
	mockedInstances    map[interface{}]*instanceMock
	provider           callMetadataProvider
	targetFunc         reflect.Value
}

func (g *mockGuard) makeCall(in []reflect.Value) []reflect.Value {
	instance, realTarget, in := g.provider(in)

	mock, found := g.mockedInstances[instance]
	if !found {
		mock, found = g.mockedInstances[nil]
	}
	if !found {
		if realTarget == nil {
			log.Fatal("Unexpected error: mocked func/method not found and real target is nil, unable to perform call")
		}
		return g.callReal(realTarget.Call, in)
	}

	var out []reflect.Value
	if mock.enabled {
		mock.RecordCall(in)

		var err error
		out, err = mock.mockedCalls.MockedOutFor(in)
		if err != nil {
			return g.defaultOut
		}
	}

	if out == nil {
		return g.callReal(mock.target.Call, in)
	}

	return out
}

func (g *mockGuard) callReal(realTarget func(in []reflect.Value) []reflect.Value, in []reflect.Value) []reflect.Value {
	g.guard.Unpatch()
	defer g.guard.Restore()
	return realTarget(in)
}

func (g *mockGuard) patchFunc(instance interface{}) (*monkey.PatchGuard, callMetadataProvider) {
	instanceType := g.targetFunc.Type()
	replacement := reflect.MakeFunc(instanceType, g.makeCall)
	mg := monkey.Patch(g.targetFunc.Interface(), replacement.Interface())

	provider := func(in []reflect.Value) (interface{}, *reflect.Value, []reflect.Value) {
		return nil, nil, in
	}

	return mg, provider
}

func (g *mockGuard) patchMethod(instance interface{}) (*monkey.PatchGuard, callMetadataProvider) {
	methodName := utils.MethodName(g.fullyQualifiedName)

	instanceType := reflect.TypeOf(instance)
	methodType, found := instanceType.MethodByName(methodName)
	if !found {
		log.Fatal("Unexpected error: the specified instance does not have a method called " + methodName)
	}

	replacement := reflect.MakeFunc(methodType.Func.Type(), g.makeCall)
	mg := monkey.PatchInstanceMethod(reflect.TypeOf(instance), methodName, replacement.Interface())

	provider := func(in []reflect.Value) (interface{}, *reflect.Value, []reflect.Value) {
		instanceValue := in[0]

		instance := instanceValue.Interface()
		realTarget := instanceValue.MethodByName(methodName)

		return instance, &realTarget, in[1:]
	}

	return mg, provider
}
