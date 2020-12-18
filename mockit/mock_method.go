package mockit

import (
	"testing"
)

// MockMethod creates a new Mock to mock an instance method
func MockMethod(t *testing.T, instance interface{}, method interface{}) Mock {
	return manager.MockMethod(t, instance, method)
}
