package mockit

import (
	"errors"
	"reflect"
	"testing"
)

func Test_findCall(t *testing.T) {
	type fields struct {
		calls      []*funcCall
		defaultOut []reflect.Value
	}
	type args struct {
		in []reflect.Value
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr error
	}{
		{
			name: "Call not found",
			fields: fields{
				calls: nil,
			},
			args: args{
				in: []reflect.Value{reflect.ValueOf("some-arg")},
			},
			want:    -1,
			wantErr: errors.New("Unable to find a call with the specified input parameters"),
		},
		{
			name: "Call found",
			fields: fields{
				calls: []*funcCall{&funcCall{
					in: []reflect.Value{reflect.ValueOf("some-arg")},
				}},
			},
			args: args{
				in: []reflect.Value{reflect.ValueOf("some-arg")},
			},
			want:    0,
			wantErr: nil,
		},
		{
			name: "Call with different number of arguments",
			fields: fields{
				calls: []*funcCall{&funcCall{
					in: []reflect.Value{reflect.ValueOf("some-arg"), reflect.ValueOf("some-other-arg")},
				}},
			},
			args: args{
				in: []reflect.Value{reflect.ValueOf("some-arg")},
			},
			want:    -1,
			wantErr: errors.New("Unable to find a call with the specified input parameters"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findCall(tt.fields.calls, tt.args.in, func(fromCalls, in []reflect.Value) bool {
				return callsMatch(fromCalls, in, false)
			})
			if err != tt.wantErr && err.Error() != tt.wantErr.Error() {
				t.Errorf("mockFunc.findCall() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("mockFunc.findCall() = %v, want %v", got, tt.want)
			}
		})
	}
}
