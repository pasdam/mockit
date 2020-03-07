package mockit

import (
	"reflect"
	"testing"
)

func Test_canBeNil(t *testing.T) {
	type args struct {
		kind reflect.Kind
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "Invalid", args: args{kind: reflect.Invalid}, want: false},
		{name: "Bool", args: args{kind: reflect.Bool}, want: false},
		{name: "Int", args: args{kind: reflect.Int}, want: false},
		{name: "Int8", args: args{kind: reflect.Int8}, want: false},
		{name: "Int16", args: args{kind: reflect.Int16}, want: false},
		{name: "Int32", args: args{kind: reflect.Int32}, want: false},
		{name: "Int64", args: args{kind: reflect.Int64}, want: false},
		{name: "Uint", args: args{kind: reflect.Uint}, want: false},
		{name: "Uint8", args: args{kind: reflect.Uint8}, want: false},
		{name: "Uint16", args: args{kind: reflect.Uint16}, want: false},
		{name: "Uint32", args: args{kind: reflect.Uint32}, want: false},
		{name: "Uint64", args: args{kind: reflect.Uint64}, want: false},
		{name: "Uintptr", args: args{kind: reflect.Uintptr}, want: false},
		{name: "Float32", args: args{kind: reflect.Float32}, want: false},
		{name: "Float64", args: args{kind: reflect.Float64}, want: false},
		{name: "Complex64", args: args{kind: reflect.Complex64}, want: false},
		{name: "Complex128", args: args{kind: reflect.Complex128}, want: false},
		{name: "Array", args: args{kind: reflect.Array}, want: true},
		{name: "Chan", args: args{kind: reflect.Chan}, want: true},
		{name: "Func", args: args{kind: reflect.Func}, want: true},
		{name: "Interface", args: args{kind: reflect.Interface}, want: true},
		{name: "Map", args: args{kind: reflect.Map}, want: true},
		{name: "Ptr", args: args{kind: reflect.Ptr}, want: true},
		{name: "Slice", args: args{kind: reflect.Slice}, want: true},
		{name: "String", args: args{kind: reflect.String}, want: false},
		{name: "Struct", args: args{kind: reflect.Struct}, want: false},
		{name: "UnsafePointer", args: args{kind: reflect.UnsafePointer}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := canBeNil(tt.args.kind); got != tt.want {
				t.Errorf("canBeNil() = %v, want %v", got, tt.want)
			}
		})
	}
}
