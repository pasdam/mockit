package argument

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAny(t *testing.T) {
	type args struct {
		arg interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Should match string",
			args: args{arg: "some string"},
		},
		{
			name: "Should match int",
			args: args{arg: 123456},
		},
		{
			name: "Should match double",
			args: args{arg: 0.123456},
		},
		{
			name: "Should match func",
			args: args{arg: TestAny},
		},
		{
			name: "Should match nil",
			args: args{arg: nil},
		},
		{
			name: "Should custom struct",
			args: args{arg: args{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Any(tt.args.arg)

			assert.True(t, got)
		})
	}
}
