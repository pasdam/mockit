package utils

import (
	"reflect"
	"testing"
)

func TestMethodFullyQualifiedName(t *testing.T) {
	type args struct {
		methodValue reflect.Value
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{
				methodValue: reflect.ValueOf(TestMethodFullyQualifiedName),
			},
			want: "github.com/pasdam/mockit/internal/utils.TestMethodFullyQualifiedName",
		},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := MethodFullyQualifiedName(tt.args.methodValue); got != tt.want {
				t.Errorf("MethodFullyQualifiedName() = %v, want %v", got, tt.want)
			}
		})
	}
}
