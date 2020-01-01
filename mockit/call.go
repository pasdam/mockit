package mockit

import (
	"reflect"
)

type call struct {
	in    []reflect.Value
	out   []reflect.Value
	count uint
}
