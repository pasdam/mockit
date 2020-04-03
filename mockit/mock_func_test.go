package mockit

import (
	"path/filepath"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/pasdam/mockit/matchers/argument"
	"github.com/stretchr/testify/assert"
)

func Test_mockFunc_ShouldRemoveTheMockWhenTheTestCompletes(t *testing.T) {
	absPath, _ := filepath.Abs(".")

	t.Run("", func(t *testing.T) {
		m := MockFunc(t, filepath.Base)
		m.With(absPath).Return("mock-value")

		assert.Equal(t, "mock-value", filepath.Base(absPath))
	})

	assert.Equal(t, "mockit", filepath.Base(absPath))
}

func Test_mockFunc_ShouldUseArgumentMatcher(t *testing.T) {
	m := MockFunc(t, filepath.Base).(*mockFunc)
	m.With(argument.Any).Return("result")

	assert.Equal(t, "result", filepath.Base("argument-1"))
	assert.Equal(t, "result", filepath.Base("argument-2"))
	assert.Equal(t, "result", filepath.Base("argument-3"))
	assert.Equal(t, 3, len(m.calls))
	m.Verify("argument-1")
	m.Verify("argument-2")
	m.Verify("argument-3")
}

func Test_mockFunc_ShouldReturnDefaultOutputIfNoMatchingCallIsFound(t *testing.T) {
	m := MockFunc(t, filepath.Base).(*mockFunc)
	m.With("argument-1").Return("out-1")
	m.With("argument-2").Return("out-2")
	m.With("argument-3").Return("out-3")

	assert.Equal(t, "", filepath.Base("non-matching-argument"))
	assert.Equal(t, 1, len(m.calls))
	m.Verify("non-matching-argument")
}

func Test_mockFunc_ShouldReturnAZeroValueIfTheMockArgumentIsNil(t *testing.T) {
	m := MockFunc(t, filepath.Walk).(*mockFunc)
	m.With("arg", nil).Return(nil)

	assert.Nil(t, filepath.Walk("arg", nil))
	assert.Equal(t, 1, len(m.calls))
	m.Verify("arg", nil)
}

func Test_mockFunc_ShouldReturnExpectedOutputIfAMatchingCallIsFound(t *testing.T) {
	m := MockFunc(t, filepath.Base).(*mockFunc)
	m.With("matching-argument").Return("some-out")

	assert.Equal(t, "some-out", filepath.Base("matching-argument"))
	assert.Equal(t, 1, len(m.calls))
	m.Verify("matching-argument")
}

func Test_mockFunc_ShouldDisableAndRestoreAMock(t *testing.T) {
	m := MockFunc(t, filepath.Base).(*mockFunc)
	m.With("matching-argument").Return("some-out")

	assert.Equal(t, "some-out", filepath.Base("matching-argument"))

	m.Disable()

	assert.Equal(t, "matching-argument", filepath.Base("matching-argument"))

	m.Enable()

	assert.Equal(t, "some-out", filepath.Base("matching-argument"))

	assert.Equal(t, 2, len(m.calls))
	m.Verify("matching-argument")
}

func Test_MockFunc(t *testing.T) {
	type args struct {
		t      *testing.T
		target interface{}
	}
	tests := []struct {
		name       string
		args       args
		defaultOut []interface{}
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
			defaultOut: []interface{}{""},
			shouldFail: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockT := new(testing.T)
			got := MockFunc(mockT, tt.args.target)
			if tt.shouldFail {
				if !mockT.Failed() {
					t.Fatalf("NewmockFunc was expected to fail, but it didn't")
				}
				if got != nil {
					t.Fatalf("NewmockFunc was expected to return nil, but it was %v", got)
				}
			} else {
				if mockT.Failed() {
					t.Fatalf("NewmockFunc wasn't expected to fail, but it did")
				}
				if got == nil {
					t.Fatalf("NewmockFunc was expected to return a valid object, but it was nil")
				}
				m := got.(*mockFunc)
				assert.Equal(t, mockT, m.t)
				assert.True(t, reflect.DeepEqual(reflect.ValueOf(tt.args.target), m.target))
				assert.Equal(t, len(tt.defaultOut), len(m.defaultOut))
				for i := 0; i < len(tt.defaultOut); i++ {
					assert.Equal(t, tt.defaultOut[i], m.defaultOut[i].Interface())
				}
			}
		})
	}
}

func Test_mockFunc_CallRealMethod(t *testing.T) {
	type fields struct {
		mocks []*funcCall
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "First mock",
			fields: fields{
				mocks: nil,
			},
		},
		{
			name: "Second mock",
			fields: fields{
				mocks: []*funcCall{&funcCall{}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &mockFunc{
				funcMockData: funcMockData{
					mocks: tt.fields.mocks,
					currentMock: &funcCall{
						in: []reflect.Value{reflect.ValueOf("some-arg")},
					},
				},
			}
			f.CallRealMethod()

			expectedMockIndex := len(tt.fields.mocks)
			assert.Equal(t, expectedMockIndex+1, len(f.mocks))
			assert.Nil(t, f.mocks[expectedMockIndex].out)
			assert.Equal(t, 1, len(f.mocks[expectedMockIndex].in))
			assert.Equal(t, "some-arg", f.mocks[expectedMockIndex].in[0].String())
			assert.Nil(t, f.calls)
			assert.Nil(t, f.defaultOut)
			assert.Nil(t, f.guard)
			assert.Nil(t, f.t)
		})
	}
}

func Test_mockFunc_Return_ShouldCallConfigureMockReturn(t *testing.T) {
	expectedMock := &mockFunc{
		funcMockData: funcMockData{
			t: t,
		},
	}
	expectedIn := []interface{}{"some-in"}
	called := false
	guard := monkey.Patch(configureMockReturn, func(f *funcMockData, in ...interface{}) {
		assert.Equal(t, &expectedMock.funcMockData, f)
		assert.Equal(t, expectedIn, in)
		called = true
	})
	defer guard.Unpatch()

	expectedMock.Return(expectedIn...)

	assert.True(t, called)
}

func Test_mockFunc_ReturnDefaults(t *testing.T) {
	type fields struct {
		mocks      []*funcCall
		defaultOut []reflect.Value
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "First mock",
			fields: fields{
				mocks: nil,
			},
		},
		{
			name: "Second mock",
			fields: fields{
				mocks: []*funcCall{&funcCall{}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &mockFunc{
				funcMockData: funcMockData{
					mocks: tt.fields.mocks,
					currentMock: &funcCall{
						in: []reflect.Value{reflect.ValueOf("some-arg")},
					},
				},
			}
			f.ReturnDefaults()

			expectedMockIndex := len(tt.fields.mocks)
			assert.Equal(t, expectedMockIndex+1, len(f.mocks))
			assert.Nil(t, f.mocks[expectedMockIndex].out)
			assert.Equal(t, 1, len(f.mocks[expectedMockIndex].in))
			assert.Equal(t, "some-arg", f.mocks[expectedMockIndex].in[0].String())
			assert.Nil(t, f.calls)
			assert.Nil(t, f.defaultOut)
			assert.Nil(t, f.guard)
			assert.Nil(t, f.t)
		})
	}
}

func Test_mockFunc_Verify_ShouldCallVerifyCallWithProvidedParameters(t *testing.T) {
	expectedMock := &mockFunc{
		funcMockData: funcMockData{
			t: t,
		},
	}
	expectedIn := []interface{}{"some-in"}
	called := false
	guard := monkey.Patch(verifyCall, func(f *funcMockData, in ...interface{}) {
		assert.Equal(t, &expectedMock.funcMockData, f)
		assert.Equal(t, expectedIn, in)
		called = true
	})
	defer guard.Unpatch()

	expectedMock.Verify(expectedIn...)

	assert.True(t, called)
}

func Test_mockFunc_With_ShouldCallConfigureMockWithAndReturnItself(t *testing.T) {
	expectedMock := &mockFunc{
		funcMockData: funcMockData{
			t: t,
		},
	}
	expectedIn := []interface{}{"some-in"}
	called := false
	guard := monkey.Patch(configureMockWith, func(f *funcMockData, values ...interface{}) {
		assert.Equal(t, &expectedMock.funcMockData, f)
		assert.Equal(t, expectedIn, values)
		called = true
	})
	defer guard.Unpatch()

	expectedMock.With(expectedIn...)

	assert.True(t, called)
}

func Test_mockFunc_makeCall(t *testing.T) {
	type fields struct {
		mocks      []*funcCall
		defaultOut []reflect.Value
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
			},
			args: args{
				in: []reflect.Value{reflect.ValueOf("some-arg")},
			},
			want: []reflect.Value{reflect.ValueOf("default-out-value")},
		},
		{
			name: "Mocked output",
			fields: fields{
				mocks: []*funcCall{&funcCall{
					in:  []reflect.Value{reflect.ValueOf("some-arg")},
					out: []reflect.Value{reflect.ValueOf("mocked-out-value")},
				}},
			},
			args: args{
				in: []reflect.Value{reflect.ValueOf("some-arg")},
			},
			want: []reflect.Value{reflect.ValueOf("mocked-out-value")},
		},
		{
			name: "Real method",
			fields: fields{
				mocks: []*funcCall{&funcCall{
					in:  []reflect.Value{reflect.ValueOf("../mockit/func_mock_test.go")},
					out: nil,
				}},
			},
			args: args{
				in: []reflect.Value{reflect.ValueOf("../mockit/func_mock_test.go")},
			},
			want: []reflect.Value{reflect.ValueOf("func_mock_test.go")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockFunc{
				funcMockData: funcMockData{
					defaultOut: tt.fields.defaultOut,
					mocks:      tt.fields.mocks,
					target:     reflect.ValueOf(filepath.Base),
				},
			}
			m.guard = monkey.Patch(makeCall, func(mock *funcMockData, in []reflect.Value, defaultOut []reflect.Value, guard *monkey.PatchGuard) []reflect.Value {
				assert.Equal(t, &m.funcMockData, mock)
				assert.True(t, callsMatch(tt.args.in, in, true))
				assert.True(t, callsMatch(tt.fields.defaultOut, defaultOut, true))
				assert.Equal(t, m.guard, guard)
				return tt.want
			})
			defer m.guard.Unpatch()

			if got := m.makeCall(tt.args.in); !callsMatch(got, tt.want, true) {
				t.Errorf("mockFunc.makeCall() = %v, want %v", got, tt.want)
			}
		})
	}
}
