package mockit

import (
	"reflect"
	"testing"

	"github.com/pasdam/mockit/matchers/argument"
)

func Test_argumentsMatch(t *testing.T) {
	mockMatcher := func(val interface{}) bool { return false }

	type args struct {
		expected      reflect.Value
		actual        reflect.Value
		enableMatcher bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Equal strings",
			args: args{
				expected:      reflect.ValueOf("some-string"),
				actual:        reflect.ValueOf("some-string"),
				enableMatcher: true,
			},
			want: true,
		},
		{
			name: "Different strings",
			args: args{
				expected:      reflect.ValueOf("some-string"),
				actual:        reflect.ValueOf("some-other-string"),
				enableMatcher: true,
			},
			want: false,
		},
		{
			name: "Equal ints",
			args: args{
				expected:      reflect.ValueOf(100),
				actual:        reflect.ValueOf(100),
				enableMatcher: true,
			},
			want: true,
		},
		{
			name: "Different ints",
			args: args{
				expected:      reflect.ValueOf(100),
				actual:        reflect.ValueOf(-100),
				enableMatcher: true,
			},
			want: false,
		},
		{
			name: "Matcher matches",
			args: args{
				expected:      reflect.ValueOf(argument.Any),
				actual:        reflect.ValueOf(-100),
				enableMatcher: true,
			},
			want: true,
		},
		{
			name: "Matcher matches, but is not ebabled",
			args: args{
				expected:      reflect.ValueOf(argument.Any),
				actual:        reflect.ValueOf(-100),
				enableMatcher: false,
			},
			want: false,
		},
		{
			name: "Matcher does not match",
			args: args{
				expected:      reflect.ValueOf(mockMatcher),
				actual:        reflect.ValueOf(-100),
				enableMatcher: true,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := argumentsMatch(tt.args.expected, tt.args.actual, tt.args.enableMatcher); got != tt.want {
				t.Errorf("argumentsMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}
