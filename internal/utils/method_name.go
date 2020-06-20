package utils

import (
	"strings"
)

// MethodName returns the method's name from it fully qualified one
func MethodName(fullyQualifiedName string) string {
	nameParts := strings.Split(fullyQualifiedName, ".")
	return strings.TrimSuffix(nameParts[len(nameParts)-1], "-fm")
}
