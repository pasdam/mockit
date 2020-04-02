package mockit

import (
	"reflect"
	"testing"
)

type funcMockData struct {
	calls    []*funcCall
	instance interface{}
	mocks    []*funcCall
	t        *testing.T
	target   reflect.Value
}
