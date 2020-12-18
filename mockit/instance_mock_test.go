package mockit

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_instanceMock_Disable(t *testing.T) {
	tests := []struct {
		name    string
		enabled bool
	}{
		{
			name:    "Should disable if enabled",
			enabled: true,
		},
		{
			name:    "Should stay disabled, if already disabled",
			enabled: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &instanceMock{
				enabled: tt.enabled,
			}

			m.Disable()

			assert.False(t, m.enabled)
		})
	}
}

func Test_instanceMock_Enable(t *testing.T) {
	tests := []struct {
		name    string
		enabled bool
	}{
		{
			name:    "Should enable if disabled",
			enabled: false,
		},
		{
			name:    "Should stay enabled, if already enabled",
			enabled: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &instanceMock{
				enabled: tt.enabled,
			}

			m.Enable()

			assert.True(t, m.enabled)
		})
	}
}

func Test_instanceMock_Verify(t *testing.T) {
	target := reflect.ValueOf(filepath.Base)
	type fields struct {
		defaultOut  []reflect.Value
		calls       [][]reflect.Value
		enabled     bool
		mockedCalls *callsIndex
		target      *reflect.Value
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
				target:     &target,
			},
			args: args{
				in: []interface{}{"some-arg"},
			},
			shouldFail: true,
		},
		{
			name: "Called with a different argument",
			fields: fields{
				calls: [][]reflect.Value{
					{reflect.ValueOf("some-arg")},
				},
				defaultOut: []reflect.Value{reflect.ValueOf("default-out-value")},
				target:     &target,
			},
			args: args{
				in: []interface{}{"some-other-arg"},
			},
			shouldFail: true,
		},
		{
			name: "Called multiple times with different arguments",
			fields: fields{
				calls: [][]reflect.Value{
					{reflect.ValueOf("some-arg-1")},
					{reflect.ValueOf("some-arg-2")},
				},
				defaultOut: []reflect.Value{reflect.ValueOf("default-out-value")},
				target:     &target,
			},
			args: args{
				in: []interface{}{"some-other-arg"},
			},
			shouldFail: true,
		},
		{
			name: "Called",
			fields: fields{
				calls: [][]reflect.Value{
					{reflect.ValueOf("some-arg")},
				},
				defaultOut: []reflect.Value{reflect.ValueOf("default-out-value")},
				target:     &target,
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
			m := &instanceMock{
				defaultOut:  tt.fields.defaultOut,
				calls:       tt.fields.calls,
				enabled:     tt.fields.enabled,
				mockedCalls: tt.fields.mockedCalls,
				t:           mockT,
				target:      tt.fields.target,
			}

			m.Verify(tt.args.in...)

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

func Test_instanceMock_With(t *testing.T) {
	type fields struct {
		target reflect.Value
	}
	type args struct {
		values []interface{}
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		want          *stubBuilder
		shouldSucceed bool
	}{
		{
			name: "Valid values",
			fields: fields{
				target: reflect.ValueOf(filepath.Base),
			},
			args: args{
				values: []interface{}{"some-value"},
			},
			want: &stubBuilder{
				args: []reflect.Value{reflect.ValueOf("some-value")},
			},
			shouldSucceed: true,
		},
		{
			name: "Invalid values: wrong number of arguments",
			fields: fields{
				target: reflect.ValueOf(filepath.HasPrefix),
			},
			args: args{
				values: []interface{}{"some-value"},
			},
			want:          nil,
			shouldSucceed: false,
		},
		{
			name: "Invalid values: wrong type of arguments",
			fields: fields{
				target: reflect.ValueOf(filepath.Base),
			},
			args: args{
				values: []interface{}{"some-value", 10},
			},
			want:          nil,
			shouldSucceed: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockT := new(testing.T)
			m := &instanceMock{
				t:      mockT,
				target: &tt.fields.target,
			}

			got := m.With(tt.args.values...).(*stubBuilder)

			if tt.shouldSucceed {
				assert.False(t, mockT.Failed())
				assert.Equal(t, m, got.mock)
				assert.False(t, got.completed)

				assert.Equal(t, len(tt.want.args), len(got.args))
				for i := 0; i < len(tt.want.args); i++ {
					assert.Equal(t, tt.want.args[i].Interface(), got.args[i].Interface())
				}

			} else {
				assert.True(t, mockT.Failed())
			}
		})
	}
}

func Test_instanceMock_RecordCall(t *testing.T) {
	type fields struct {
		calls [][]reflect.Value
	}
	type args struct {
		in []reflect.Value
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   [][]reflect.Value
	}{
		{
			name: "First mocked call",
			fields: fields{
				calls: nil,
			},
			args: args{
				in: []reflect.Value{reflect.ValueOf("some-first-value")},
			},
			want: [][]reflect.Value{
				{reflect.ValueOf("some-first-value")},
			},
		},
		{
			name: "Second mocked call",
			fields: fields{
				calls: [][]reflect.Value{
					{reflect.ValueOf("some-first-value")},
				},
			},
			args: args{
				in: []reflect.Value{reflect.ValueOf("some-second-value")},
			},
			want: [][]reflect.Value{
				{reflect.ValueOf("some-first-value")},
				{reflect.ValueOf("some-second-value")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &instanceMock{
				calls: tt.fields.calls,
			}

			m.RecordCall(tt.args.in)

			assert.Equal(t, len(tt.want), len(m.calls))
			for i := 0; i < len(tt.want); i++ {
				assert.Equal(t, len(tt.want[i]), len(m.calls[i]))

				for j := 0; j < len(tt.want[i]); j++ {
					assert.Equal(t, tt.want[i][j].Interface(), m.calls[i][j].Interface())
				}
			}
		})
	}
}
