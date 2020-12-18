package mockit

import (
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

func emptyProvider(guard *mockGuard) func(instance interface{}) (*monkey.PatchGuard, callMetadataProvider) {
	return func(instance interface{}) (*monkey.PatchGuard, callMetadataProvider) {
		return nil, nil
	}
}

func Test_mockManager_mock_shouldFailTestIfTargetFunctionIsNil(t *testing.T) {
	manager := mockManager{
		mockedTypes: make(map[string]*mockGuard),
	}
	mockT := new(testing.T)
	instance := "some-instance"

	manager.mock(mockT, instance, nil, emptyProvider)

	assert.True(t, mockT.Failed())
}

func Test_mockManager_mock_shouldFailTestIfTargetIsNotAFunction(t *testing.T) {
	manager := mockManager{
		mockedTypes: make(map[string]*mockGuard),
	}
	mockT := new(testing.T)
	instance := "some-instance"

	manager.mock(mockT, instance, "some-non-func-target", emptyProvider)

	assert.True(t, mockT.Failed())
}
