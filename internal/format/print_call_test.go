package format_test

import (
	"reflect"
	"testing"

	"github.com/pasdam/mockit/internal/format"
	"github.com/stretchr/testify/assert"
)

func Test_printCall(t *testing.T) {
	type t1 struct {
		a string
		b int
		c bool
		d float64
	}
	type args struct {
		target reflect.Value
		in     []reflect.Value
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Nil arguments",
			args: args{
				target: reflect.ValueOf(Test_printCall),
				in:     nil,
			},
			want: "Test_printCall()",
		},
		{
			name: "Empty arguments",
			args: args{
				target: reflect.ValueOf(format.PrintCall),
				in:     []reflect.Value{},
			},
			want: "PrintCall()",
		},
		{
			name: "With 1 argument",
			args: args{
				target: reflect.ValueOf(format.PrintCall),
				in:     []reflect.Value{reflect.ValueOf("some-arg")},
			},
			want: "PrintCall(some-arg)",
		},
		{
			name: "With 2 arguments",
			args: args{
				target: reflect.ValueOf(format.PrintCall),
				in:     []reflect.Value{reflect.ValueOf("some-arg-1"), reflect.ValueOf("some-arg-2")},
			},
			want: "PrintCall(some-arg-1, some-arg-2)",
		},
		{
			name: "With struct argument",
			args: args{
				target: reflect.ValueOf(format.PrintCall),
				in: []reflect.Value{reflect.ValueOf(t1{
					a: "some-a-val",
					b: 123,
					c: true,
					d: 45.6,
				})},
			},
			want: "PrintCall({a:some-a-val b:123 c:true d:45.6})",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := format.PrintCall(&tt.args.target, tt.args.in)

			assert.Equal(t, tt.want, got)
		})
	}
}
