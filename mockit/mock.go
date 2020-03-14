package mockit

import (
	"testing"
)

// Mock contains methods to mock a call with specified arguments, and verify it
type Mock interface {

	// Disable disable the mock, so interactions will be with real objects
	Disable()

	// Enable restore the mock
	Enable()

	// Verify fails the test if a call with the specified arguments wasn't made
	Verify(t *testing.T, in ...interface{})

	// With configures the mock to respond to the specified arguments
	With(values ...interface{}) Stub
}
