package mockit

import (
	"errors"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/pasdam/mockit/matchers/argument"
	"github.com/stretchr/testify/assert"
)

func Test_MockFunc_ShouldUseArgumentMatcher(t *testing.T) {
	m := NewMockFunc(t, filepath.Base).(*mockFunc)
	m.Mock(t, []interface{}{argument.Any}, []interface{}{"result"})

	assert.Equal(t, "result", filepath.Base("argument-1"))
	assert.Equal(t, "result", filepath.Base("argument-2"))
	assert.Equal(t, "result", filepath.Base("argument-3"))
	assert.Equal(t, 3, len(m.calls))
	m.Verify(t, []interface{}{"argument-1"})
	m.Verify(t, []interface{}{"argument-2"})
	m.Verify(t, []interface{}{"argument-3"})
}

func Test_MockFunc_ShouldReturnDefaultOutputIfNoMatchingCallIsFound(t *testing.T) {
	m := NewMockFunc(t, filepath.Base).(*mockFunc)
	m.Mock(t, []interface{}{"argument-1"}, []interface{}{"out-1"})
	m.Mock(t, []interface{}{"argument-2"}, []interface{}{"out-2"})
	m.Mock(t, []interface{}{"argument-3"}, []interface{}{"out-3"})

	assert.Equal(t, "", filepath.Base("non-matching-argument"))
	assert.Equal(t, 1, len(m.calls))
	m.Verify(t, []interface{}{"non-matching-argument"})
}

func Test_MockFunc_ShouldReturnAZeroValueIfTheMockArgumentIsNil(t *testing.T) {
	m := NewMockFunc(t, filepath.Walk).(*mockFunc)
	m.Mock(t, []interface{}{"arg", nil}, []interface{}{nil})

	assert.Nil(t, filepath.Walk("arg", nil))
	assert.Equal(t, 1, len(m.calls))
	m.Verify(t, []interface{}{"arg", nil})
}

func Test_MockFunc_ShouldReturnExpectedOutputIfAMatchingCallIsFound(t *testing.T) {
	m := NewMockFunc(t, filepath.Base).(*mockFunc)
	m.Mock(t, []interface{}{"matching-argument"}, []interface{}{"some-out"})

	assert.Equal(t, "some-out", filepath.Base("matching-argument"))
	assert.Equal(t, 1, len(m.calls))
	m.Verify(t, []interface{}{"matching-argument"})
}

func Test_MockFunc_ShouldDisableAndRestoreAMock(t *testing.T) {
	m := NewMockFunc(t, filepath.Base).(*mockFunc)
	m.Mock(t, []interface{}{"matching-argument"}, []interface{}{"some-out"})

	assert.Equal(t, "some-out", filepath.Base("matching-argument"))

	m.UnMock()

	assert.Equal(t, "matching-argument", filepath.Base("matching-argument"))

	m.Restore()

	assert.Equal(t, "some-out", filepath.Base("matching-argument"))

	assert.Equal(t, 2, len(m.calls))
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
		mocks      []*call
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
		expectedMocks int
		shouldFail    bool
	}{
		{
			name: "Invalid in",
			fields: fields{
				mocks:      nil,
				defaultOut: []reflect.Value{},
			},
			args: args{
				in:  []interface{}{"some-in", "some-additional-in"},
				out: []interface{}{"some-out"},
			},
			expectedMocks: 0,
			shouldFail:    true,
		},
		{
			name: "Invalid out",
			fields: fields{
				mocks:      nil,
				defaultOut: []reflect.Value{},
			},
			args: args{
				in:  []interface{}{"some-in"},
				out: []interface{}{"some-out", "some-additional-out"},
			},
			expectedMocks: 0,
			shouldFail:    true,
		},
		{
			name: "Valid mock call",
			fields: fields{
				mocks:      nil,
				defaultOut: []reflect.Value{},
			},
			args: args{
				in:  []interface{}{"some-in"},
				out: []interface{}{"some-out"},
			},
			expectedMocks: 1,
			shouldFail:    false,
		},
		{
			name: "Override existing call",
			fields: fields{
				mocks: []*call{
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
			expectedMocks: 1,
			shouldFail:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockT := new(testing.T)
			m := &mockFunc{
				mocks:      tt.fields.mocks,
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
			assert.Equal(t, tt.expectedMocks, len(m.mocks))
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
			name: "Not called",
			fields: fields{
				calls:      nil,
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
					in:  []reflect.Value{reflect.ValueOf("some-arg")},
					out: []reflect.Value{reflect.ValueOf("mocked-out-value")},
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
				calls:      nil,
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
			got, err := findCall(m.calls, tt.args.in, func(fromCalls, in []reflect.Value) bool {
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

func Test_mockFunc_makeCall(t *testing.T) {
	type fields struct {
		mocks      []*call
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
				mocks:      nil,
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
				mocks: []*call{&call{
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
				mocks:      tt.fields.mocks,
				defaultOut: tt.fields.defaultOut,
				typeOf:     tt.fields.typeOf,
			}
			if got := m.makeCall(tt.args.in); !callsMatch(got, tt.want, true) {
				t.Errorf("mockFunc.makeCall() = %v, want %v", got, tt.want)
			}
			if len(m.calls) != 1 {
				t.Errorf("Expected 1 mocked call, got %d", len(m.calls))
			}
			// if m.calls[0].count != 1 {
			// 	t.Errorf("Expected 1 recorded call, got %d", len(m.calls))
			// }
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
				calls:      nil,
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
			m.recordCall(tt.args.in)
			expectedCallsCount := len(tt.fields.calls) + 1
			assert.Equal(t, expectedCallsCount, len(m.calls))
			assert.Equal(t, tt.args.in, m.calls[expectedCallsCount-1].in)
			assert.Nil(t, m.calls[expectedCallsCount-1].out)
			// if m.calls[0].count != 1 {
			// 	t.Errorf("Expected 1 recorded call, got %d", len(m.calls))
			// }
		})
	}
}
