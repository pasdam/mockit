package mockit

import (
	"reflect"
)

type funcCall struct {
	in  []reflect.Value
	out []reflect.Value
}
