package mockit

import (
	"testing"
)

// MockFunc creates a new Mock to mock a function
func MockFunc(t *testing.T, targetFn interface{}) Mock {
	return manager.MockFunc(t, targetFn)
}
