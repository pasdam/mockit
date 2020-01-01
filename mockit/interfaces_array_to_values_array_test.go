package mockit

import (
	"reflect"
	"testing"
)

func Test_interfacesArrayToValuesArray(t *testing.T) {
	type args struct {
		args []interface{}
	}
	tests := []struct {
		name string
		args args
		want []reflect.Value
	}{
		{
			name: "String",
			args: args{
				args: []interface{}{"some-arg"},
			},
			want: []reflect.Value{reflect.ValueOf("some-arg")},
		},
		{
			name: "Int",
			args: args{
				args: []interface{}{1234},
			},
			want: []reflect.Value{reflect.ValueOf(1234)},
		},
		{
			name: "Multiple arguments",
			args: args{
				args: []interface{}{"some-arg", 1234},
			},
			want: []reflect.Value{reflect.ValueOf("some-arg"), reflect.ValueOf(1234)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := interfacesArrayToValuesArray(tt.args.args)
			if len(got) != len(tt.want) {
				t.Errorf("Expected result length (%d) is different than the actual one (%d)", len(got), len(tt.want))
			}
			for i := 0; i < len(tt.want); i++ {
				if tt.want[i].Interface() != got[i].Interface() {
					t.Errorf("interfaceArrayToValue() = %v, want %v", got, tt.want)
					return
				}
			}
		})
	}
}
