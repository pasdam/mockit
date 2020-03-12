package mockit

import (
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

func Test_callsMatch(t *testing.T) {
	type args struct {
		expected       []reflect.Value
		actual         []reflect.Value
		enableMatchers bool
	}
	tests := []struct {
		name          string
		args          args
		argumentMatch []bool
		want          bool
	}{
		{
			name: "Expected length lower than actual length",
			args: args{
				expected:       []reflect.Value{reflect.ValueOf("a"), reflect.ValueOf("b")},
				actual:         []reflect.Value{reflect.ValueOf("a"), reflect.ValueOf("b"), reflect.ValueOf("c")},
				enableMatchers: true,
			},
			want: false,
		},
		{
			name: "Expected length greater than actual length",
			args: args{
				expected:       []reflect.Value{reflect.ValueOf("a"), reflect.ValueOf("b"), reflect.ValueOf("c")},
				actual:         []reflect.Value{reflect.ValueOf("a"), reflect.ValueOf("b")},
				enableMatchers: true,
			},
			want: false,
		},
		{
			name: "First argument does not match",
			args: args{
				expected:       []reflect.Value{reflect.ValueOf("a"), reflect.ValueOf("b"), reflect.ValueOf("c")},
				actual:         []reflect.Value{reflect.ValueOf("not-matching"), reflect.ValueOf("b"), reflect.ValueOf("c")},
				enableMatchers: true,
			},
			argumentMatch: []bool{false},
			want:          false,
		},
		{
			name: "Second argument does not match",
			args: args{
				expected:       []reflect.Value{reflect.ValueOf("a"), reflect.ValueOf("b"), reflect.ValueOf("c")},
				actual:         []reflect.Value{reflect.ValueOf("a"), reflect.ValueOf("not-matching"), reflect.ValueOf("c")},
				enableMatchers: true,
			},
			argumentMatch: []bool{true, false},
			want:          false,
		},
		{
			name: "Argument matches",
			args: args{
				expected:       []reflect.Value{reflect.ValueOf("a"), reflect.ValueOf("b"), reflect.ValueOf("c")},
				actual:         []reflect.Value{reflect.ValueOf("a"), reflect.ValueOf("b"), reflect.ValueOf("c")},
				enableMatchers: true,
			},
			argumentMatch: []bool{true, true, true},
			want:          true,
		},
	}
	for _, tt := range tests {
		callsCount := 0
		monkey.Patch(argumentsMatch, func(expected reflect.Value, actual reflect.Value, enableMatchers bool) bool {
			assert.Equal(t, tt.args.expected[callsCount].Interface(), expected.Interface())
			assert.Equal(t, tt.args.actual[callsCount].Interface(), actual.Interface())
			assert.Equal(t, tt.args.enableMatchers, enableMatchers)

			result := tt.argumentMatch[callsCount]
			callsCount = callsCount + 1
			return result
		})
		defer monkey.Unpatch(argumentsMatch)

		t.Run(tt.name, func(t *testing.T) {
			if got := callsMatch(tt.args.expected, tt.args.actual, tt.args.enableMatchers); got != tt.want {
				t.Errorf("callsMatch() = %v, want %v", got, tt.want)
			}
			assert.Equal(t, len(tt.argumentMatch), callsCount)
		})
	}
}
