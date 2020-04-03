package mockit

import (
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

type mockMethodGuardTest struct{}

func (m *mockMethodGuardTest) TestMethod(arg int) string { return "real-method-value" }

func Test_newMockMethodGuard(t *testing.T) {
	mockInstance := &mockMethodGuardTest{}
	type args struct {
		fullyQualifiedName string
		method             reflect.Value
		instance           interface{}
	}
	type fields struct {
		defaultOut []reflect.Value
	}
	tests := []struct {
		name    string
		args    args
		fields  fields
		wantErr bool
	}{
		{
			name: "Without error",
			args: args{
				fullyQualifiedName: "github.com/pasdam/mockit/mockit.(*mockMethodGuardTest).TestMethod-fm",
				method:             reflect.ValueOf(mockInstance.TestMethod),
				instance:           mockInstance,
			},
			fields: fields{
				defaultOut: []reflect.Value{reflect.ValueOf("")},
			},
			wantErr: false,
		},
		{
			name: "With error",
			args: args{
				fullyQualifiedName: "not-existing-method",
				method:             reflect.ValueOf(mockInstance.TestMethod),
				instance:           mockInstance,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockT := new(testing.T)

			got := newMockMethodGuard(mockT, tt.args.fullyQualifiedName, tt.args.method, tt.args.instance)

			assert.Equal(t, tt.wantErr, mockT.Failed())
			if !tt.wantErr {
				assert.NotNil(t, got)
				assert.Equal(t, len(got.defaultOut), len(tt.fields.defaultOut))
				for i := 0; i < len(got.defaultOut); i++ {
					assert.Equal(t, got.defaultOut[i].Interface(), tt.fields.defaultOut[i].Interface())
				}
				assert.NotNil(t, got.guard)
				assert.NotNil(t, got.mocks)
				assert.Empty(t, got.mocks)
				assert.Equal(mockT, t, got.t)

			} else {
				assert.Nil(t, got)
			}
		})
	}
}

func Test_mockMethodGuard_methodMock(t *testing.T) {
	defaultOut := []reflect.Value{reflect.ValueOf("some-out")}
	instance := &mockMethodGuardTest{}
	mockT := new(testing.T)
	type fields struct {
		mock       *mockMethod
		defaultOut []reflect.Value
	}
	tests := []struct {
		name   string
		fields fields
		want   *mockMethod
	}{
		{
			name: "Existing mock",
			fields: fields{
				mock: &mockMethod{
					funcMockData: funcMockData{
						t: t,
					},
				},
			},
			want: &mockMethod{
				funcMockData: funcMockData{
					t: t,
				},
			},
		},
		{
			name: "Create mock",
			fields: fields{
				mock:       nil,
				defaultOut: defaultOut,
			},
			want: &mockMethod{
				enabled: true,
				funcMockData: funcMockData{
					instance:   instance,
					t:          mockT,
					defaultOut: defaultOut,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockMethodGuard{
				defaultOut: tt.fields.defaultOut,
				mocks:      make(map[interface{}]*mockMethod),
				t:          mockT,
			}
			if tt.fields.mock != nil {
				m.mocks[instance] = tt.fields.mock
			}

			got := m.methodMock(instance)

			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_mockMethodGuard_makeCall(t *testing.T) {
	instance := &mockMethodGuardTest{}
	type fields struct {
		defaultOut []reflect.Value
		methodName string
		mock       *mockMethod
	}
	type args struct {
		in []reflect.Value
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []reflect.Value
	}{
		{
			name: "Instance not mocked, it should call real method",
			fields: fields{
				defaultOut: nil,
				methodName: "TestMethod",
				mock:       nil,
			},
			args: args{
				in: []reflect.Value{reflect.ValueOf(10)},
			},
			want: []reflect.Value{reflect.ValueOf("real-method-value")},
		},
		{
			name: "Instance mocked but the mock is disabled, it should call real method",
			fields: fields{
				defaultOut: []reflect.Value{reflect.ValueOf("default-out")},
				methodName: "TestMethod",
				mock: &mockMethod{
					enabled: false,
					funcMockData: funcMockData{
						mocks:  nil,
						target: reflect.ValueOf(instance.TestMethod),
					},
				},
			},
			args: args{
				in: []reflect.Value{reflect.ValueOf(-50)},
			},
			want: []reflect.Value{reflect.ValueOf("real-method-value")},
		},
		{
			name: "Instance mocked, it should return mocked value",
			fields: fields{
				defaultOut: []reflect.Value{reflect.ValueOf("default-out")},
				mock: &mockMethod{
					enabled: true,
					funcMockData: funcMockData{
						mocks:  nil,
						target: reflect.ValueOf(instance.TestMethod),
					},
				},
			},
			args: args{
				in: []reflect.Value{reflect.ValueOf(-50)},
			},
			want: []reflect.Value{reflect.ValueOf("mocked-makecall-value")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGuard := &mockMethodGuard{
				defaultOut: tt.fields.defaultOut,
				methodName: tt.fields.methodName,
				mocks:      make(map[interface{}]*mockMethod),
				t:          t,
			}
			mockGuard.guard = monkey.PatchInstanceMethod(reflect.TypeOf(instance), "TestMethod", func(*mockMethodGuardTest, int) string {
				return "mocked-method-value"
			})
			defer mockGuard.guard.Unpatch()
			if tt.fields.mock != nil {
				mockGuard.mocks[instance] = tt.fields.mock
			}
			inWithInstance := []reflect.Value{reflect.ValueOf(instance)}
			inWithInstance = append(inWithInstance, tt.args.in...)
			guardMC := monkey.Patch(makeCall, func(actualMock *funcMockData, in []reflect.Value, defaultOut []reflect.Value, guard *monkey.PatchGuard) []reflect.Value {
				assert.Equal(t, &tt.fields.mock.funcMockData, actualMock)
				assert.True(t, callsMatch(tt.args.in, in, true))
				assert.True(t, callsMatch(tt.fields.defaultOut, defaultOut, true))
				assert.Equal(t, mockGuard.guard, guard)
				return []reflect.Value{reflect.ValueOf("mocked-makecall-value")}
			})
			defer guardMC.Unpatch()

			if got := mockGuard.makeCall(inWithInstance); !callsMatch(got, tt.want, true) {
				t.Errorf("mockMethodGuard.makeCall() = %v, want %v", got, tt.want)
			}
		})
	}
}
