package mockit

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_configureMockReturn(t *testing.T) {
	type args struct {
		values []interface{}
	}
	type fields struct {
		mocks []*funcCall
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		shouldFail bool
	}{
		{
			name: "First mock",
			fields: fields{
				mocks: nil,
			},
			args: args{
				values: []interface{}{"out-1"},
			},
			shouldFail: false,
		},
		{
			name: "Second mock",
			fields: fields{
				mocks: []*funcCall{&funcCall{}},
			},
			args: args{
				values: []interface{}{"out-2"},
			},
			shouldFail: false,
		},
		{
			name: "Wrong return type",
			fields: fields{
				mocks: []*funcCall{&funcCall{}},
			},
			args: args{
				values: []interface{}{100},
			},
			shouldFail: true,
		},
		{
			name: "Not enough return values",
			fields: fields{
				mocks: []*funcCall{&funcCall{}},
			},
			args: args{
				values: []interface{}{},
			},
			shouldFail: true,
		},
		{
			name: "Too many return values",
			fields: fields{
				mocks: []*funcCall{&funcCall{}},
			},
			args: args{
				values: []interface{}{"out-0", "out-1"},
			},
			shouldFail: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := []reflect.Value{reflect.ValueOf("some-arg")}
			mockT := new(testing.T)
			f := &funcMockData{
				mocks:  tt.fields.mocks,
				target: reflect.ValueOf(filepath.Base),
				t:      mockT,
				currentMock: &funcCall{
					in: in,
				},
			}
			configureMockReturn(f, tt.args.values...)

			assert.Equal(t, tt.shouldFail, mockT.Failed())

			if !tt.shouldFail {
				expectedMockIndex := len(tt.fields.mocks)
				assert.Equal(t, expectedMockIndex+1, len(f.mocks))
				assert.Equal(t, in, f.mocks[expectedMockIndex].in)

				assert.Equal(t, 1, len(f.mocks[expectedMockIndex].out))
				assert.Equal(t, tt.args.values[0], f.mocks[expectedMockIndex].out[0].String())

				assert.Nil(t, f.calls)
			}
		})
	}
}
