package mockit

import (
	"reflect"
)

func defaultFuncOutput(typeOf reflect.Type) []reflect.Value {
	out := make([]reflect.Value, 0, typeOf.NumOut())
	for i := 0; i < typeOf.NumOut(); i++ {
		out = append(out, reflect.Zero(typeOf.Out(i)))
	}
	return out
}
