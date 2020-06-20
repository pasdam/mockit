package mockit

import (
	"reflect"
	"strings"

	"github.com/pasdam/mockit/internal/format"
)

func verifyCall(f *funcMockData, in ...interface{}) {
	inValues := interfacesArrayToValuesArray(in, f.target.Type().In)
	_, err := findCall(f.calls, inValues, func(fromCalls, in []reflect.Value) bool {
		return callsMatch(in, fromCalls, true)
	})
	if err != nil {
		builder := strings.Builder{}
		builder.WriteString("Expected call: ")
		builder.WriteString(format.PrintCall(f.target, inValues))
		if len(f.calls) > 0 {
			builder.WriteString("; but it recorded the following instead: ")
			builder.WriteString(format.PrintCall(f.target, f.calls[0].in))
			for i := 1; i < len(f.calls); i++ {
				builder.WriteString(", ")
				builder.WriteString(format.PrintCall(f.target, f.calls[i].in))
			}

		} else {
			builder.WriteString("; but no call was recorded")
		}
		f.t.Error(builder.String())
	}
}
