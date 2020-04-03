package mockit

import (
	"strings"
)

func getMethodName(fullyQualifiedName string) string {
	nameParts := strings.Split(fullyQualifiedName, ".")
	return strings.TrimSuffix(nameParts[len(nameParts)-1], "-fm")
}
