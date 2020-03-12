package mockit

import (
	"reflect"
	"testing"
)

func voidFunction() {
}

func functionWith1Out() string {
	return "some-out"
}

func functionWith4Out() (int, float64, bool, reflect.Value) {
	return 12345, 0.12345, true, reflect.Value{}
}

func Test_defaultFuncOutput(t *testing.T) {
	var zeroValue reflect.Value
	type args struct {
		typeOf reflect.Type
	}
	tests := []struct {
		name string
		args args
		want []reflect.Value
	}{
		{
			name: "Void func",
			args: args{
				typeOf: reflect.TypeOf(voidFunction),
			},
			want: []reflect.Value{},
		},
		{
			name: "1 output",
			args: args{
				typeOf: reflect.TypeOf(functionWith1Out),
			},
			want: []reflect.Value{reflect.ValueOf("")},
		},
		{
			name: "4 outputs",
			args: args{
				typeOf: reflect.TypeOf(functionWith4Out),
			},
			want: []reflect.Value{reflect.ValueOf(0), reflect.ValueOf(0.0), reflect.ValueOf(false), reflect.ValueOf(zeroValue)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := defaultFuncOutput(tt.args.typeOf); !callsMatch(got, tt.want, true) {
				t.Errorf("defaultFuncOutput() = %v, want %v", got, tt.want)
			}
		})
	}
}
