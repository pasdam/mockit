package mockit

import (
	"errors"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

func convertToValuesAndVerifiesTestFunc01(a int, b float64, c string) {}

func Test_convertToValuesAndVerifies(t *testing.T) {
	f01 := reflect.TypeOf(convertToValuesAndVerifiesTestFunc01)

	type args struct {
		values                []interface{}
		expectedValuesCount   int
		expectedValueProvider func(int) reflect.Type
	}
	tests := []struct {
		name    string
		args    args
		want    []reflect.Value
		wantErr error
	}{
		{
			name: "Convert all values",
			args: args{
				values:                []interface{}{0, 1.2, "some-string"},
				expectedValuesCount:   f01.NumIn(),
				expectedValueProvider: f01.In,
			},
			want: []reflect.Value{
				reflect.ValueOf(0),
				reflect.ValueOf(1.2),
				reflect.ValueOf("some-string"),
			},
			wantErr: nil,
		},
		{
			name: "Should fail if verifyValues returns error",
			args: args{
				values:                []interface{}{120, 341.2, "some-other-string"},
				expectedValuesCount:   f01.NumIn(),
				expectedValueProvider: f01.In,
			},
			want:    nil,
			wantErr: errors.New("some-verify-error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockT := new(testing.T)

			if tt.wantErr != nil {
				guard := monkey.Patch(verifyValues, func(a int, b func(int) reflect.Type, c []reflect.Value) error {
					return tt.wantErr
				})
				t.Cleanup(guard.Unpatch)
			}

			got := convertToValuesAndVerifies(mockT, tt.args.values, tt.args.expectedValuesCount, tt.args.expectedValueProvider)

			if tt.want != nil {
				for i := 0; i < len(tt.args.values); i++ {
					assert.Equal(t, tt.want[i].Interface(), got[i].Interface())
				}
			} else {
				assert.Nil(t, got)
			}
		})
	}
}
