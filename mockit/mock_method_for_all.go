package mockit

import (
	"testing"
)

// MockMethodForAll creates a new Mock to mock the method for any instance of
// the specified type
func MockMethodForAll(t *testing.T, instance interface{}, method interface{}) Mock {
	return manager.MockMethodForAll(t, instance, method)
}
