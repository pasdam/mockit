package mockit

import (
	"reflect"
	"testing"
)

func Test_inputsMatch(t *testing.T) {
	type args struct {
		expected []reflect.Value
		actual   []reflect.Value
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Nil arrays",
			args: args{
				expected: nil,
				actual:   nil,
			},
			want: true,
		},
		{
			name: "Different length",
			args: args{
				expected: []reflect.Value{reflect.ValueOf("some-content"), reflect.ValueOf("some-other-content")},
				actual:   []reflect.Value{reflect.ValueOf("some-content")},
			},
			want: false,
		},
		{
			name: "Different content",
			args: args{
				expected: []reflect.Value{reflect.ValueOf("some-content")},
				actual:   []reflect.Value{reflect.ValueOf("some-other-content")},
			},
			want: false,
		},
		{
			name: "Same length and content",
			args: args{
				expected: []reflect.Value{reflect.ValueOf("some-content")},
				actual:   []reflect.Value{reflect.ValueOf("some-content")},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := valuesArrayMatch(tt.args.expected, tt.args.actual); got != tt.want {
				t.Errorf("inputsMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}
