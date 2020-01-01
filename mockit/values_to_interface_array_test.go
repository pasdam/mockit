package mockit

import (
	"reflect"
	"testing"
)

func Test_valuesToInterfaceArray(t *testing.T) {
	type args struct {
		args []reflect.Value
	}
	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		{
			name: "String",
			args: args{
				args: []reflect.Value{reflect.ValueOf("some-arg")},
			},
			want: []interface{}{"some-arg"},
		},
		{
			name: "Int",
			args: args{
				args: []reflect.Value{reflect.ValueOf(1234)},
			},
			want: []interface{}{1234},
		},
		{
			name: "Multiple arguments",
			args: args{
				args: []reflect.Value{reflect.ValueOf("some-arg"), reflect.ValueOf(1234)},
			},
			want: []interface{}{"some-arg", 1234},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := valuesToInterfaceArray(tt.args.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("valuesToInterfaceArray() = %v, want %v", got, tt.want)
			}
		})
	}
}
