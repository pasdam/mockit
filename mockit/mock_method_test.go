package mockit

import (
	"errors"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

type mockMethodTest struct{}

func (m *mockMethodTest) TestMethod(arg int) string { return "real-method-value" }

func Test_MockMethod_ShouldReturnExpectedValueForTheReadmeExample(t *testing.T) {
	// NOTE: if this fails (i.e. the contract changed), please update the README as well
	err := errors.New("some-error")
	m := MockMethod(t, err, err.Error)
	m.With().Return("some-other-value")
	assert.Equal(t, "some-other-value", err.Error())
}

func Test_MockMethod_ShouldMockOnlyTheSpecifiedInstance(t *testing.T) {
	err1 := errors.New("some-error")
	m := MockMethod(t, err1, err1.Error)
	m.With().Return("some-other-value")
	assert.Equal(t, "some-other-value", err1.Error())

	err2 := errors.New("some-other-error")
	assert.Equal(t, "some-other-error", err2.Error())
}

func TestMockMethod(t *testing.T) {
	instance := &mockMethodTest{}
	emg := &mockMethodGuard{
		mocks: make(map[interface{}]*mockMethod),
	}
	em := &mockMethod{
		funcMockData: funcMockData{
			t: t,
		},
	}
	target := instance.TestMethod
	emg.mocks[instance] = em
	type fields struct {
		fullyQualifiedName string
		existingMock       *mockMethodGuard
	}
	type args struct {
		t        *testing.T
		instance interface{}
		method   interface{}
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		want       Mock
		shouldFail bool
	}{
		{
			name: "Should fail if instance is nil",
			fields: fields{
				fullyQualifiedName: "",
				existingMock:       nil,
			},
			args: args{
				instance: nil,
				method:   target,
			},
			want:       nil,
			shouldFail: true,
		},
		{
			name: "Should fail if method is nil",
			fields: fields{
				fullyQualifiedName: "",
				existingMock:       nil,
			},
			args: args{
				instance: instance,
				method:   nil,
			},
			want:       nil,
			shouldFail: true,
		},
		{
			name: "Should fail if method is not a func",
			fields: fields{
				fullyQualifiedName: "",
				existingMock:       nil,
			},
			args: args{
				instance: instance,
				method:   instance,
			},
			want:       nil,
			shouldFail: true,
		},
		{
			name: "Should return existing mock",
			fields: fields{
				fullyQualifiedName: "github.com/pasdam/mockit/mockit.(*mockMethodTest).TestMethod-fm",
				existingMock:       emg,
			},
			args: args{
				instance: instance,
				method:   target,
			},
			want:       em,
			shouldFail: false,
		},
		{
			name: "Should create new mock if instance is not mocked",
			fields: fields{
				fullyQualifiedName: "github.com/pasdam/mockit/mockit.(*mockMethodTest).TestMethod-fm",
				existingMock:       nil,
			},
			args: args{
				instance: instance,
				method:   target,
			},
			want: &mockMethod{
				enabled: true,
				funcMockData: funcMockData{
					instance:   instance,
					target:     reflect.ValueOf(target),
					defaultOut: []reflect.Value{reflect.ValueOf("")},
				},
			},
			shouldFail: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockT := new(testing.T)
			methodMocksMap = make(map[string]*mockMethodGuard)
			if tt.fields.existingMock != nil {
				methodMocksMap[tt.fields.fullyQualifiedName] = tt.fields.existingMock
			}
			if tt.want != nil {
				tt.want.(*mockMethod).t = mockT
			}

			got := MockMethod(mockT, tt.args.instance, tt.args.method)

			assert.Equal(t, tt.shouldFail, mockT.Failed())

			if tt.want != nil {
				gotWithoutDefaultOut := *got.(*mockMethod)
				gotWithoutDefaultOut.defaultOut = nil
				wantWithoutDefaultOut := *tt.want.(*mockMethod)
				wantWithoutDefaultOut.defaultOut = nil

				assert.Equal(t, &wantWithoutDefaultOut, &gotWithoutDefaultOut)
				assert.True(t, callsMatch(tt.want.(*mockMethod).defaultOut, got.(*mockMethod).defaultOut, true))

			} else {
				assert.Nil(t, got)
			}
		})
	}
}

func Test_newMethodMock(t *testing.T) {
	instance := &mockMethodTest{}

	got := newMethodMock(t, instance)

	assert.Equal(t, &mockMethod{
		enabled: true,
		funcMockData: funcMockData{
			instance: instance,
			t:        t,
		},
	}, got)
}

func Test_methodMock_ShouldDisableMockAtTheEndOfTheTest(t *testing.T) {
	var m *mockMethod

	t.Run("", func(t *testing.T) {
		m = newMethodMock(t, t)
		assert.True(t, m.enabled)
	})

	assert.False(t, m.enabled)
}

func Test_mockMethod_CallRealMethod(t *testing.T) {
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
				mocks: []*funcCall{{}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &mockMethod{
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
			assert.Nil(t, f.t)
		})
	}
}

func Test_mockMethod_Disable(t *testing.T) {
	actual := mockMethod{enabled: true}

	actual.Disable()

	assert.Equal(t, mockMethod{enabled: false}, actual)
}

func Test_mockMethod_Enable(t *testing.T) {
	actual := mockMethod{enabled: false}

	actual.Enable()

	assert.Equal(t, mockMethod{enabled: true}, actual)
}

func Test_mockMethod_Return_ShouldCallConfigureMockReturn(t *testing.T) {
	expectedMock := &mockMethod{
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

func Test_mockMethod_ReturnDefaults(t *testing.T) {
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
				mocks: []*funcCall{{}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &mockMethod{
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
			assert.Nil(t, f.t)
		})
	}
}

func Test_mockMethod_Verify_ShouldCallVerifyCallWithProvidedParameters(t *testing.T) {
	expectedMock := &mockMethod{
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

func Test_mockMethod_With_ShouldCallConfigureMockWithAndReturnItself(t *testing.T) {
	expectedMock := &mockMethod{
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
