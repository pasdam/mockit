package format

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pasdam/mockit/internal/utils"
)

// PrintCall prints the specified method call
func PrintCall(target *reflect.Value, in []reflect.Value) string {
	var str strings.Builder

	str.WriteString(utils.MethodName(utils.MethodFullyQualifiedName(*target)))
	str.WriteString("(")
	if len(in) > 0 {
		str.WriteString(fmt.Sprintf("%+v", in[0]))
		for i := 1; i < len(in); i++ {
			str.WriteString(", ")
			str.WriteString(fmt.Sprintf("%+v", in[i]))
		}
	}
	str.WriteString(")")

	return str.String()
}
