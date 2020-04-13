package argument

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCaptor_Capture(t *testing.T) {
	type fields struct {
		Value interface{}
	}
	type args struct {
		arg interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Should capture string value",
			fields: fields{
				Value: nil,
			},
			args: args{
				arg: "some-value",
			},
		},
		{
			name: "Should capture int value",
			fields: fields{
				Value: nil,
			},
			args: args{
				arg: 123456,
			},
		},
		{
			name: "Should capture double value",
			fields: fields{
				Value: nil,
			},
			args: args{
				arg: 0.123456,
			},
		},
		{
			name: "Should capture value and override existing one",
			fields: fields{
				Value: "some-old-value",
			},
			args: args{
				arg: "some-overriding-value",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Captor{
				Value: tt.fields.Value,
			}

			got := c.Capture(tt.args.arg)

			assert.True(t, got)
			assert.Equal(t, tt.args.arg, c.Value)
		})
	}
}
