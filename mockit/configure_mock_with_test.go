package mockit

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_configureMockWith(t *testing.T) {
	type args struct {
		in []interface{}
	}
	tests := []struct {
		name       string
		args       args
		shouldFail bool
	}{
		{
			name: "Invalid arguments",
			args: args{
				in: []interface{}{"some-in", "some-additional-in"},
			},
			shouldFail: true,
		},
		{
			name: "Valid arguments",
			args: args{
				in: []interface{}{"some-in"},
			},
			shouldFail: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockT := new(testing.T)
			m := &funcMockData{
				target: reflect.ValueOf(filepath.Base),
				t:      mockT,
			}
			configureMockWith(m, tt.args.in...)

			if tt.shouldFail != mockT.Failed() {
				if tt.shouldFail {
					t.Errorf("Verify was expected to fail, but it didn't")
				} else {
					t.Errorf("Verify wasn't expected to fail, but it did")
				}
			}
			assert.Equal(t, 0, len(m.mocks))
			assert.NotNil(t, m.currentMock)
			if tt.shouldFail == false {
				assert.Equal(t, len(tt.args.in), len(m.currentMock.in))
				for i := 0; i < len(tt.args.in); i++ {
					assert.Equal(t, tt.args.in[i], m.currentMock.in[i].Interface())
				}
			} else {
				assert.Nil(t, m.currentMock.in)
			}
			assert.Nil(t, m.currentMock.out)
		})
	}
}
