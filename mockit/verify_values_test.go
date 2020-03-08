package mockit

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/pasdam/mockit/matchers/argument"
)

func Test_verify_values(t *testing.T) {
	type args struct {
		expectedCount         int
		expectedValueProvider func(int) reflect.Type
		actualValues          []reflect.Value
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "Different sizes",
			args: args{
				expectedCount:         2,
				expectedValueProvider: nil,
				actualValues:          []reflect.Value{reflect.ValueOf("some-value")},
			},
			wantErr: errors.New("Expected values count (2) is different than the actual size (1)"),
		},
		{
			name: "Different type",
			args: args{
				expectedCount:         2,
				expectedValueProvider: reflect.TypeOf(os.Setenv).In,
				actualValues:          []reflect.Value{reflect.ValueOf("some-value"), reflect.ValueOf(12345)},
			},
			wantErr: errors.New("Type at index 1 is different than expected (string): actual type int"),
		},
		{
			name: "Different types, but assignable",
			args: args{
				expectedCount:         1,
				expectedValueProvider: reflect.TypeOf(os.Setenv).Out,
				actualValues:          []reflect.Value{reflect.ValueOf(errors.New("some-error"))},
			},
			wantErr: nil,
		},
		{
			name: "Same types",
			args: args{
				expectedCount:         2,
				expectedValueProvider: reflect.TypeOf(os.Setenv).In,
				actualValues:          []reflect.Value{reflect.ValueOf("some-value"), reflect.ValueOf("some-other-value")},
			},
			wantErr: nil,
		},
		{
			name: "Interface type, nil value",
			args: args{
				expectedCount:         1,
				expectedValueProvider: reflect.TypeOf(os.Setenv).Out,
				actualValues:          []reflect.Value{reflect.Value{}},
			},
			wantErr: nil,
		},
		{
			name: "Pointer type, nil value",
			args: args{
				expectedCount:         1,
				expectedValueProvider: reflect.TypeOf(os.Environ).Out,
				actualValues:          []reflect.Value{reflect.Value{}},
			},
			wantErr: nil,
		},
		{
			name: "Nil value for a non interface/pointer type",
			args: args{
				expectedCount:         1,
				expectedValueProvider: reflect.TypeOf(os.Getpid).Out,
				actualValues:          []reflect.Value{reflect.Value{}},
			},
			wantErr: errors.New("Cannot assign nil at index 0 to the type int"),
		},
		{
			name: "With matcher",
			args: args{
				expectedCount:         1,
				expectedValueProvider: reflect.TypeOf(os.Getpid).Out,
				actualValues:          []reflect.Value{reflect.ValueOf(argument.Any)},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := verifyValues(tt.args.expectedCount, tt.args.expectedValueProvider, tt.args.actualValues)
			if tt.wantErr == nil && err != nil || tt.wantErr != nil && (err == nil || err != nil && tt.wantErr.Error() != err.Error()) {
				t.Errorf("verify_values() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
