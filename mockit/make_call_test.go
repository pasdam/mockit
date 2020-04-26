package mockit

import (
	"path/filepath"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

func Test_makeCall(t *testing.T) {
	type args struct {
		mock       *funcMockData
		in         []reflect.Value
		defaultOut []reflect.Value
	}
	tests := []struct {
		name      string
		args      args
		want      []reflect.Value
		wantCount int
	}{
		{
			name: "Default output",
			args: args{
				mock: &funcMockData{
					mocks: nil,
				},
				in:         []reflect.Value{reflect.ValueOf("some-arg")},
				defaultOut: []reflect.Value{reflect.ValueOf("default-out-value")},
			},
			want: []reflect.Value{reflect.ValueOf("default-out-value")},
		},
		{
			name: "Mocked output",
			args: args{
				mock: &funcMockData{
					mocks: []*funcCall{{
						in:  []reflect.Value{reflect.ValueOf("some-arg")},
						out: []reflect.Value{reflect.ValueOf("mocked-out-value")},
					}},
				},
				in: []reflect.Value{reflect.ValueOf("some-arg")},
			},
			want: []reflect.Value{reflect.ValueOf("mocked-out-value")},
		},
		{
			name: "Real method",
			args: args{
				mock: &funcMockData{
					mocks: []*funcCall{{
						in:  []reflect.Value{reflect.ValueOf("../mockit/func_mock_test.go")},
						out: nil,
					}},
				},
				in: []reflect.Value{reflect.ValueOf("../mockit/func_mock_test.go")},
			},
			want: []reflect.Value{reflect.ValueOf("func_mock_test.go")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			guard := monkey.Patch(filepath.Base, func(string) string { return "wrong-result" })
			defer guard.Unpatch()

			mock := &funcMockData{
				mocks:  tt.args.mock.mocks,
				target: reflect.ValueOf(filepath.Base),
			}

			got := makeCall(mock, tt.args.in, tt.args.defaultOut, guard)

			assert.True(t, callsMatch(got, tt.want, true))
			assert.Equal(t, 1, len(mock.calls))
			assert.Equal(t, tt.args.in, mock.calls[0].in)
			assert.Nil(t, mock.calls[0].out)
		})
	}
}

func Test_makeCall_shouldRecordMultipleCalls(t *testing.T) {
	m := &funcMockData{}
	defaultOut := []reflect.Value{reflect.ValueOf("")}
	guard := monkey.Patch(filepath.Base, func(string) string { return "" })
	performCall := func(arg string) {
		makeCall(m, []reflect.Value{reflect.ValueOf(arg)}, defaultOut, guard)
	}

	performCall("arg-0")
	performCall("arg-1")
	performCall("arg-2")

	assert.Equal(t, 3, len(m.calls))

	assert.Equal(t, 1, len(m.calls[0].in))
	assert.Equal(t, "arg-0", m.calls[0].in[0].String())
	assert.Nil(t, m.calls[0].out)

	assert.Equal(t, 1, len(m.calls[1].in))
	assert.Equal(t, "arg-1", m.calls[1].in[0].String())
	assert.Nil(t, m.calls[1].out)

	assert.Equal(t, 1, len(m.calls[2].in))
	assert.Equal(t, "arg-2", m.calls[2].in[0].String())
	assert.Nil(t, m.calls[2].out)
}
