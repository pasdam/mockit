package utils_test

import (
	"testing"

	"github.com/pasdam/mockit/internal/utils"
	"github.com/stretchr/testify/assert"
)

func Test_MethodName(t *testing.T) {
	type args struct {
		fullyQualifiedName string
	}
	tests := []struct {
		name               string
		args               args
		fullyQualifiedName string
		want               string
	}{
		{
			name: "With package and suffix",
			args: args{
				fullyQualifiedName: "reflect.Value.Addr-fm",
			},
			want: "Addr",
		},
		{
			name: "With package but no suffix",
			args: args{
				fullyQualifiedName: "reflect.Value.Addr",
			},
			want: "Addr",
		},
		{
			name: "Without package and with suffix",
			args: args{
				fullyQualifiedName: "Addr-fm",
			},
			want: "Addr",
		},
		{
			name: "Without package and suffix",
			args: args{
				fullyQualifiedName: "Addr",
			},
			want: "Addr",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			methodName := utils.MethodName(tt.args.fullyQualifiedName)

			assert.Equal(t, tt.want, methodName)
		})
	}
}
