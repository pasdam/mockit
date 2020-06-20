package mockit

import (
	"path/filepath"
	"reflect"
	"testing"
)

func Test_verifyCall(t *testing.T) {
	type fields struct {
		calls      []*funcCall
		defaultOut []reflect.Value
		target     reflect.Value
	}
	type args struct {
		in []interface{}
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		shouldFail bool
	}{
		{
			name: "Not called",
			fields: fields{
				calls:      nil,
				defaultOut: []reflect.Value{reflect.ValueOf("default-out-value")},
				target:     reflect.ValueOf(filepath.Base),
			},
			args: args{
				in: []interface{}{"some-arg"},
			},
			shouldFail: true,
		},
		{
			name: "Called with a different argument",
			fields: fields{
				calls: []*funcCall{{
					in:  []reflect.Value{reflect.ValueOf("some-arg")},
					out: []reflect.Value{reflect.ValueOf("mocked-out-value")},
				}},
				defaultOut: []reflect.Value{reflect.ValueOf("default-out-value")},
				target:     reflect.ValueOf(filepath.Base),
			},
			args: args{
				in: []interface{}{"some-other-arg"},
			},
			shouldFail: true,
		},
		{
			name: "Called multiple times with different arguments",
			fields: fields{
				calls: []*funcCall{
					{
						in:  []reflect.Value{reflect.ValueOf("some-arg-1")},
						out: []reflect.Value{reflect.ValueOf("mocked-out-value-1")},
					},
					{
						in:  []reflect.Value{reflect.ValueOf("some-arg-2")},
						out: []reflect.Value{reflect.ValueOf("mocked-out-value-2")},
					},
				},
				defaultOut: []reflect.Value{reflect.ValueOf("default-out-value")},
				target:     reflect.ValueOf(filepath.Base),
			},
			args: args{
				in: []interface{}{"some-other-arg"},
			},
			shouldFail: true,
		},
		{
			name: "Called",
			fields: fields{
				calls: []*funcCall{{
					in:  []reflect.Value{reflect.ValueOf("some-arg")},
					out: []reflect.Value{reflect.ValueOf("mocked-out-value")},
				}},
				defaultOut: []reflect.Value{reflect.ValueOf("default-out-value")},
				target:     reflect.ValueOf(filepath.Base),
			},
			args: args{
				in: []interface{}{"some-arg"},
			},
			shouldFail: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockT := new(testing.T)
			m := &funcMockData{
				calls:  tt.fields.calls,
				target: tt.fields.target,
				t:      mockT,
			}

			verifyCall(m, tt.args.in...)

			// TODO: verify log message
			if mockT.Failed() != tt.shouldFail {
				if tt.shouldFail {
					t.Errorf("Verify was expected to fail, but it didn't")
				} else {
					t.Errorf("Verify wasn't expected to fail, but it did")
				}
			}
		})
	}
}
