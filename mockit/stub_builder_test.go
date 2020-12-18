package mockit

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func stubBuilderReturnTestFunc01() (string, int) { return "", 0 }

func Test_stubBuilder_CallRealMethod(t *testing.T) {
	target := reflect.ValueOf(os.Setenv)
	type fields struct {
		args      []reflect.Value
		mock      *instanceMock
		completed bool
	}
	tests := []struct {
		name          string
		fields        fields
		shouldSucceed bool
		wantMocks     *callsIndex
	}{
		{
			name: "Should fail if the stubbing is already completed",
			fields: fields{
				args: []reflect.Value{reflect.ValueOf("some-value")},
				mock: &instanceMock{
					defaultOut:  []reflect.Value{},
					calls:       [][]reflect.Value{},
					enabled:     true,
					mockedCalls: &callsIndex{},
					target:      &target,
				},
				completed: true,
			},
			shouldSucceed: false,
			wantMocks:     nil,
		},
		{
			name: "Should complete the stub",
			fields: fields{
				args: []reflect.Value{reflect.ValueOf("some-value")},
				mock: &instanceMock{
					defaultOut:  nil,
					calls:       nil,
					enabled:     false,
					mockedCalls: &callsIndex{},
					target:      nil,
				},
				completed: false,
			},
			shouldSucceed: true,
			wantMocks: &callsIndex{
				in: [][]reflect.Value{
					{reflect.ValueOf("some-value")},
				},
				out: [][]reflect.Value{
					nil,
				},
			},
		},
		{
			name: "Should complete the stub, with existing mocked calls",
			fields: fields{
				args: []reflect.Value{reflect.ValueOf("some-value")},
				mock: &instanceMock{
					defaultOut: nil,
					calls:      nil,
					enabled:    false,
					mockedCalls: &callsIndex{
						in: [][]reflect.Value{
							{reflect.ValueOf("some-previous-input")},
						},
						out: [][]reflect.Value{
							{reflect.ValueOf("some-previous-output")},
						},
					},
					target: nil,
				},
				completed: false,
			},
			shouldSucceed: true,
			wantMocks: &callsIndex{
				in: [][]reflect.Value{
					{reflect.ValueOf("some-previous-input")},
					{reflect.ValueOf("some-value")},
				},
				out: [][]reflect.Value{
					{reflect.ValueOf("some-previous-output")},
					nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &stubBuilder{
				args:      tt.fields.args,
				mock:      tt.fields.mock,
				completed: tt.fields.completed,
			}
			b.mock.t = new(testing.T)
			b.CallRealMethod()
			if tt.shouldSucceed {
				assert.False(t, b.mock.t.Failed())

				assert.Equal(t, len(tt.wantMocks.in), len(b.mock.mockedCalls.in))
				for i := 0; i < len(b.mock.mockedCalls.in); i++ {
					for j := 0; j < len(b.mock.mockedCalls.in[i]); j++ {
						assert.Equal(t, tt.wantMocks.in[i][j].Interface(), b.mock.mockedCalls.in[i][j].Interface())
					}
				}

				assert.Equal(t, len(tt.wantMocks.out), len(b.mock.mockedCalls.out))
				for i := 0; i < len(b.mock.mockedCalls.out); i++ {
					for j := 0; j < len(b.mock.mockedCalls.out[i]); j++ {
						assert.Equal(t, tt.wantMocks.out[i][j].Interface(), b.mock.mockedCalls.out[i][j].Interface())
					}
				}

			} else {
				assert.True(t, b.mock.t.Failed())
			}
		})
	}
}

func Test_stubBuilder_Return(t *testing.T) {
	target := reflect.ValueOf(stubBuilderReturnTestFunc01)
	type fields struct {
		args      []reflect.Value
		mock      *instanceMock
		completed bool
	}
	type args struct {
		values []interface{}
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		shouldSucceed bool
		wantMocks     *callsIndex
	}{
		{
			name: "Should fail if the stubbing is already completed",
			fields: fields{
				args: []reflect.Value{reflect.ValueOf("some-value")},
				mock: &instanceMock{
					defaultOut:  []reflect.Value{},
					calls:       [][]reflect.Value{},
					enabled:     true,
					mockedCalls: &callsIndex{},
					target:      &target,
				},
				completed: true,
			},
			args: args{
				values: []interface{}{"some-fail-value"},
			},
			shouldSucceed: false,
			wantMocks:     nil,
		},
		{
			name: "Should complete the stub",
			fields: fields{
				args: []reflect.Value{reflect.ValueOf("some-value")},
				mock: &instanceMock{
					defaultOut:  nil,
					calls:       nil,
					enabled:     false,
					mockedCalls: &callsIndex{},
					target:      &target,
				},
				completed: false,
			},
			args: args{
				values: []interface{}{"some-success-value", 100},
			},
			shouldSucceed: true,
			wantMocks: &callsIndex{
				in: [][]reflect.Value{
					{reflect.ValueOf("some-value")},
				},
				out: [][]reflect.Value{
					{reflect.ValueOf("some-success-value"), reflect.ValueOf(100)},
				},
			},
		},
		{
			name: "Should complete the stub, with existing mocked calls",
			fields: fields{
				args: []reflect.Value{reflect.ValueOf("some-value")},
				mock: &instanceMock{
					defaultOut: nil,
					calls:      nil,
					enabled:    false,
					mockedCalls: &callsIndex{
						in: [][]reflect.Value{
							{reflect.ValueOf("some-previous-input")},
						},
						out: [][]reflect.Value{
							{reflect.ValueOf("some-previous-output")},
						},
					},
					target: &target,
				},
				completed: false,
			},
			args: args{
				values: []interface{}{"some-other-success-value", 200},
			},
			shouldSucceed: true,
			wantMocks: &callsIndex{
				in: [][]reflect.Value{
					{reflect.ValueOf("some-previous-input")},
					{reflect.ValueOf("some-value")},
				},
				out: [][]reflect.Value{
					{reflect.ValueOf("some-previous-output")},
					{reflect.ValueOf("some-other-success-value"), reflect.ValueOf(200)},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &stubBuilder{
				args:      tt.fields.args,
				mock:      tt.fields.mock,
				completed: tt.fields.completed,
			}
			b.mock.t = new(testing.T)
			b.Return(tt.args.values...)
			if tt.shouldSucceed {
				assert.False(t, b.mock.t.Failed())

				assert.Equal(t, len(tt.wantMocks.in), len(b.mock.mockedCalls.in))
				for i := 0; i < len(b.mock.mockedCalls.in); i++ {
					for j := 0; j < len(b.mock.mockedCalls.in[i]); j++ {
						assert.Equal(t, tt.wantMocks.in[i][j].Interface(), b.mock.mockedCalls.in[i][j].Interface())
					}
				}

				assert.Equal(t, len(tt.wantMocks.out), len(b.mock.mockedCalls.out))
				for i := 0; i < len(b.mock.mockedCalls.out); i++ {
					for j := 0; j < len(b.mock.mockedCalls.out[i]); j++ {
						assert.Equal(t, tt.wantMocks.out[i][j].Interface(), b.mock.mockedCalls.out[i][j].Interface())
					}
				}

			} else {
				assert.True(t, b.mock.t.Failed())
			}
		})
	}
}

func Test_stubBuilder_ReturnDefaults(t *testing.T) {
	target := reflect.ValueOf(stubBuilderReturnTestFunc01)
	type fields struct {
		args      []reflect.Value
		mock      *instanceMock
		completed bool
	}
	tests := []struct {
		name          string
		fields        fields
		shouldSucceed bool
		wantMocks     *callsIndex
	}{
		{
			name: "Should fail if the stubbing is already completed",
			fields: fields{
				args: []reflect.Value{reflect.ValueOf("some-value")},
				mock: &instanceMock{
					defaultOut:  []reflect.Value{},
					calls:       [][]reflect.Value{},
					enabled:     true,
					mockedCalls: &callsIndex{},
					target:      &target,
				},
				completed: true,
			},
			shouldSucceed: false,
			wantMocks:     nil,
		},
		{
			name: "Should complete the stub",
			fields: fields{
				args: []reflect.Value{reflect.ValueOf("some-value")},
				mock: &instanceMock{
					defaultOut:  []reflect.Value{reflect.ValueOf("some-success-value"), reflect.ValueOf(300)},
					calls:       nil,
					enabled:     false,
					mockedCalls: &callsIndex{},
					target:      nil,
				},
				completed: false,
			},
			shouldSucceed: true,
			wantMocks: &callsIndex{
				in: [][]reflect.Value{
					{reflect.ValueOf("some-value")},
				},
				out: [][]reflect.Value{
					{reflect.ValueOf("some-success-value"), reflect.ValueOf(300)},
				},
			},
		},
		{
			name: "Should complete the stub, with existing mocked calls",
			fields: fields{
				args: []reflect.Value{reflect.ValueOf("some-value")},
				mock: &instanceMock{
					defaultOut: []reflect.Value{reflect.ValueOf("some-other-success-value"), reflect.ValueOf(400)},
					calls:      nil,
					enabled:    false,
					mockedCalls: &callsIndex{
						in: [][]reflect.Value{
							{reflect.ValueOf("some-previous-input")},
						},
						out: [][]reflect.Value{
							{reflect.ValueOf("some-previous-output")},
						},
					},
					target: nil,
				},
				completed: false,
			},
			shouldSucceed: true,
			wantMocks: &callsIndex{
				in: [][]reflect.Value{
					{reflect.ValueOf("some-previous-input")},
					{reflect.ValueOf("some-value")},
				},
				out: [][]reflect.Value{
					{reflect.ValueOf("some-previous-output")},
					{reflect.ValueOf("some-other-success-value"), reflect.ValueOf(400)},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &stubBuilder{
				args:      tt.fields.args,
				mock:      tt.fields.mock,
				completed: tt.fields.completed,
			}
			b.mock.t = new(testing.T)
			b.ReturnDefaults()
			if tt.shouldSucceed {
				assert.False(t, b.mock.t.Failed())

				assert.Equal(t, len(tt.wantMocks.in), len(b.mock.mockedCalls.in))
				for i := 0; i < len(b.mock.mockedCalls.in); i++ {
					for j := 0; j < len(b.mock.mockedCalls.in[i]); j++ {
						assert.Equal(t, tt.wantMocks.in[i][j].Interface(), b.mock.mockedCalls.in[i][j].Interface())
					}
				}

				assert.Equal(t, len(tt.wantMocks.out), len(b.mock.mockedCalls.out))
				for i := 0; i < len(b.mock.mockedCalls.out); i++ {
					for j := 0; j < len(b.mock.mockedCalls.out[i]); j++ {
						assert.Equal(t, tt.wantMocks.out[i][j].Interface(), b.mock.mockedCalls.out[i][j].Interface())
					}
				}

			} else {
				assert.True(t, b.mock.t.Failed())
			}
		})
	}
}
