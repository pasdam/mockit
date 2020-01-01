package mockit

import (
	"errors"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MockFunc_ShouldReturnDefaultOutputIfNoMatchingCallIsFound(t *testing.T) {
	m := NewMockFunc(t, filepath.Base).(*mockFunc)
	m.Mock(t, []interface{}{"argument-1"}, []interface{}{"out-1"})
	m.Mock(t, []interface{}{"argument-2"}, []interface{}{"out-2"})
	m.Mock(t, []interface{}{"argument-3"}, []interface{}{"out-3"})

	assert.Equal(t, "", filepath.Base("non-matching-argument"))
	assert.Equal(t, 4, len(m.calls))
	m.Verify(t, []interface{}{"non-matching-argument"})
}

func Test_MockFunc_ShouldReturnExpectedOutputIfAMatchingCallIsFound(t *testing.T) {
	m := NewMockFunc(t, filepath.Base).(*mockFunc)
	m.Mock(t, []interface{}{"matching-argument"}, []interface{}{"some-out"})

	assert.Equal(t, "some-out", filepath.Base("matching-argument"))
	assert.Equal(t, 1, len(m.calls))
	assert.Equal(t, 1, int(m.calls[0].count))
	m.Verify(t, []interface{}{"matching-argument"})
}

func Test_NewMockFunc(t *testing.T) {
	type args struct {
		t      *testing.T
		target interface{}
	}
	tests := []struct {
		name       string
		args       args
		shouldFail bool
	}{
		{
			name: "Non function",
			args: args{
				target: "non-function-type",
			},
			shouldFail: true,
		},
		{
			name: "Function",
			args: args{
				target: filepath.Base,
			},
			shouldFail: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockT := new(testing.T)
			got := NewMockFunc(mockT, tt.args.target)
			if tt.shouldFail {
				if !mockT.Failed() {
					t.Errorf("NewMockFunc was expected to fail, but it didn't")
				}
				if got != nil {
					t.Errorf("NewMockFunc was expected to return nil, but it was %v", got)
				}
			} else {
				if mockT.Failed() {
					t.Errorf("NewMockFunc wasn't expected to fail, but it did")
				}
				if got == nil {
					t.Errorf("NewMockFunc was expected to return a valid object, but it was nil")
				}
			}
		})
	}
}

func Test_mockFunc_Mock(t *testing.T) {
	type fields struct {
		calls      []*call
		defaultOut []reflect.Value
	}
	type args struct {
		in  []interface{}
		out []interface{}
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		expectedCalls int
		shouldFail    bool
	}{
		{
			name: "Invalid in",
			fields: fields{
				calls:      []*call{},
				defaultOut: []reflect.Value{},
			},
			args: args{
				in:  []interface{}{"some-in", "some-additional-in"},
				out: []interface{}{"some-out"},
			},
			expectedCalls: 0,
			shouldFail:    true,
		},
		{
			name: "Invalid out",
			fields: fields{
				calls:      []*call{},
				defaultOut: []reflect.Value{},
			},
			args: args{
				in:  []interface{}{"some-in"},
				out: []interface{}{"some-out", "some-additional-out"},
			},
			expectedCalls: 0,
			shouldFail:    true,
		},
		{
			name: "Valid mock call",
			fields: fields{
				calls:      []*call{},
				defaultOut: []reflect.Value{},
			},
			args: args{
				in:  []interface{}{"some-in"},
				out: []interface{}{"some-out"},
			},
			expectedCalls: 1,
			shouldFail:    false,
		},
		{
			name: "Override existing call",
			fields: fields{
				calls: []*call{
					&call{
						in:  []reflect.Value{reflect.ValueOf("some-in")},
						out: []reflect.Value{reflect.ValueOf("some-old-out")},
					},
				},
				defaultOut: []reflect.Value{},
			},
			args: args{
				in:  []interface{}{"some-in"},
				out: []interface{}{"some-out"},
			},
			expectedCalls: 1,
			shouldFail:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockT := new(testing.T)
			m := &mockFunc{
				calls:      tt.fields.calls,
				defaultOut: tt.fields.defaultOut,
				typeOf:     reflect.TypeOf(filepath.Base),
			}
			m.Mock(mockT, tt.args.in, tt.args.out)

			if tt.shouldFail != mockT.Failed() {
				if tt.shouldFail {
					t.Errorf("Verify was expected to fail, but it didn't")
				} else {
					t.Errorf("Verify wasn't expected to fail, but it did")
				}
			}
			if tt.expectedCalls != len(m.calls) {
				t.Errorf("Expected number of calls %d is different than actual (%d)", tt.expectedCalls, len(m.calls))
			}
		})
	}
}

func Test_mockFunc_Verify(t *testing.T) {
	type fields struct {
		calls      []*call
		defaultOut []reflect.Value
		typeOf     reflect.Type
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
			name: "Not mocked",
			fields: fields{
				calls:      []*call{},
				defaultOut: []reflect.Value{reflect.ValueOf("default-out-value")},
				typeOf:     reflect.TypeOf(filepath.Base),
			},
			args: args{
				in: []interface{}{"some-arg"},
			},
			shouldFail: true,
		},
		{
			name: "Not called",
			fields: fields{
				calls: []*call{&call{
					in:    []reflect.Value{reflect.ValueOf("some-arg")},
					out:   []reflect.Value{reflect.ValueOf("mocked-out-value")},
					count: 0,
				}},
				defaultOut: []reflect.Value{reflect.ValueOf("default-out-value")},
				typeOf:     reflect.TypeOf(filepath.Base),
			},
			args: args{
				in: []interface{}{"some-arg"},
			},
			shouldFail: true,
		},
		{
			name: "Called with a different argument",
			fields: fields{
				calls: []*call{&call{
					in:  []reflect.Value{reflect.ValueOf("some-arg")},
					out: []reflect.Value{reflect.ValueOf("mocked-out-value")},
				}},
				defaultOut: []reflect.Value{reflect.ValueOf("default-out-value")},
				typeOf:     reflect.TypeOf(filepath.Base),
			},
			args: args{
				in: []interface{}{"some-other-arg"},
			},
			shouldFail: true,
		},
		{
			name: "Called",
			fields: fields{
				calls: []*call{&call{
					in:    []reflect.Value{reflect.ValueOf("some-arg")},
					out:   []reflect.Value{reflect.ValueOf("mocked-out-value")},
					count: 1,
				}},
				defaultOut: []reflect.Value{reflect.ValueOf("default-out-value")},
				typeOf:     reflect.TypeOf(filepath.Base),
			},
			args: args{
				in: []interface{}{"some-arg"},
			},
			shouldFail: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockFunc{
				calls:      tt.fields.calls,
				defaultOut: tt.fields.defaultOut,
				typeOf:     tt.fields.typeOf,
			}
			mockT := new(testing.T)
			m.Verify(mockT, tt.args.in)
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

func Test_mockFunc_findCall(t *testing.T) {
	type fields struct {
		calls      []*call
		defaultOut []reflect.Value
		typeOf     reflect.Type
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
				calls:      []*call{},
				defaultOut: []reflect.Value{reflect.ValueOf("default-out-value")},
				typeOf:     reflect.TypeOf(filepath.Base),
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
				calls: []*call{&call{
					in:  []reflect.Value{reflect.ValueOf("some-arg")},
					out: []reflect.Value{reflect.ValueOf("mocked-out-value")},
				}},
				defaultOut: []reflect.Value{reflect.ValueOf("default-out-value")},
				typeOf:     reflect.TypeOf(filepath.Base),
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
				calls: []*call{&call{
					in:  []reflect.Value{reflect.ValueOf("some-arg"), reflect.ValueOf("some-other-arg")},
					out: []reflect.Value{reflect.ValueOf("mocked-out-value")},
				}},
				defaultOut: []reflect.Value{reflect.ValueOf("default-out-value")},
				typeOf:     reflect.TypeOf(filepath.Base),
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
			m := &mockFunc{
				calls:      tt.fields.calls,
				defaultOut: tt.fields.defaultOut,
				typeOf:     tt.fields.typeOf,
			}
			got, err := m.findCall(tt.args.in)
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

func Test_mockFunc_makeCall(t *testing.T) {
	type fields struct {
		calls      []*call
		defaultOut []reflect.Value
		typeOf     reflect.Type
	}
	type args struct {
		in []reflect.Value
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      []reflect.Value
		wantCount int
	}{
		{
			name: "Default output",
			fields: fields{
				calls:      []*call{},
				defaultOut: []reflect.Value{reflect.ValueOf("default-out-value")},
				typeOf:     reflect.TypeOf(filepath.Base),
			},
			args: args{
				in: []reflect.Value{reflect.ValueOf("some-arg")},
			},
			want: []reflect.Value{reflect.ValueOf("default-out-value")},
		},
		{
			name: "Mocked output",
			fields: fields{
				calls: []*call{&call{
					in:  []reflect.Value{reflect.ValueOf("some-arg")},
					out: []reflect.Value{reflect.ValueOf("mocked-out-value")},
				}},
				defaultOut: []reflect.Value{reflect.ValueOf("default-out-value")},
				typeOf:     reflect.TypeOf(filepath.Base),
			},
			args: args{
				in: []reflect.Value{reflect.ValueOf("some-arg")},
			},
			want: []reflect.Value{reflect.ValueOf("mocked-out-value")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockFunc{
				calls:      tt.fields.calls,
				defaultOut: tt.fields.defaultOut,
				typeOf:     tt.fields.typeOf,
			}
			if got := m.makeCall(tt.args.in); !valuesArrayMatch(got, tt.want) {
				t.Errorf("mockFunc.makeCall() = %v, want %v", got, tt.want)
			}
			if len(m.calls) != 1 {
				t.Errorf("Expected 1 mocked call, got %d", len(m.calls))
			}
			if m.calls[0].count != 1 {
				t.Errorf("Expected 1 recorded call, got %d", len(m.calls))
			}
		})
	}
}

func Test_mockFunc_recordCall(t *testing.T) {
	type fields struct {
		calls      []*call
		defaultOut []reflect.Value
		typeOf     reflect.Type
	}
	type args struct {
		in []reflect.Value
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Default output",
			fields: fields{
				calls:      []*call{},
				defaultOut: []reflect.Value{reflect.ValueOf("default-out-value")},
				typeOf:     reflect.TypeOf(filepath.Base),
			},
			args: args{
				in: []reflect.Value{reflect.ValueOf("some-arg")},
			},
		},
		{
			name: "Mocked output",
			fields: fields{
				calls: []*call{&call{
					in:  []reflect.Value{reflect.ValueOf("some-arg")},
					out: []reflect.Value{reflect.ValueOf("mocked-out-value")},
				}},
				defaultOut: []reflect.Value{reflect.ValueOf("default-out-value")},
				typeOf:     reflect.TypeOf(filepath.Base),
			},
			args: args{
				in: []reflect.Value{reflect.ValueOf("some-arg")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockFunc{
				calls:      tt.fields.calls,
				defaultOut: tt.fields.defaultOut,
				typeOf:     tt.fields.typeOf,
			}
			got := m.recordCall(tt.args.in)
			if len(m.calls) != 1 {
				t.Errorf("Expected 1 mocked call, got %d", len(m.calls))
			}
			if m.calls[0] != got {
				t.Errorf("mockFunc.recordCall() = %v, want %v", got, m.calls[0])
			}
			if m.calls[0].count != 1 {
				t.Errorf("Expected 1 recorded call, got %d", len(m.calls))
			}
		})
	}
}
