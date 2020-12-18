package mockit

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_callsIndex_Add(t *testing.T) {
	type fields struct {
		in  [][]reflect.Value
		out [][]reflect.Value
	}
	type args struct {
		in  []reflect.Value
		out []reflect.Value
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   fields
	}{
		{
			name: "Fist entry",
			fields: fields{
				in:  nil,
				out: nil,
			},
			args: args{
				in:  []reflect.Value{reflect.ValueOf("some-first-in-value"), reflect.ValueOf(100)},
				out: []reflect.Value{reflect.ValueOf("some-first-out-value"), reflect.ValueOf(200)},
			},
			want: fields{
				in: [][]reflect.Value{
					{reflect.ValueOf("some-first-in-value"), reflect.ValueOf(100)},
				},
				out: [][]reflect.Value{
					{reflect.ValueOf("some-first-out-value"), reflect.ValueOf(200)},
				},
			},
		},
		{
			name: "Second entry",
			fields: fields{
				in: [][]reflect.Value{
					{reflect.ValueOf("some-first-in-value"), reflect.ValueOf(100)},
				},
				out: [][]reflect.Value{
					{reflect.ValueOf("some-first-out-value"), reflect.ValueOf(200)},
				},
			},
			args: args{
				in:  []reflect.Value{reflect.ValueOf("some-second-in-value"), reflect.ValueOf(300)},
				out: []reflect.Value{reflect.ValueOf("some-second-out-value"), reflect.ValueOf(400)},
			},
			want: fields{
				in: [][]reflect.Value{
					{reflect.ValueOf("some-first-in-value"), reflect.ValueOf(100)},
					{reflect.ValueOf("some-second-in-value"), reflect.ValueOf(300)},
				},
				out: [][]reflect.Value{
					{reflect.ValueOf("some-first-out-value"), reflect.ValueOf(200)},
					{reflect.ValueOf("some-second-out-value"), reflect.ValueOf(400)},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := &callsIndex{
				in:  tt.fields.in,
				out: tt.fields.out,
			}
			idx.Add(tt.args.in, tt.args.out)

			assert.Equal(t, len(tt.want.in), len(idx.in))
			for i := 0; i < len(tt.want.in); i++ {
				assert.Equal(t, len(tt.want.in[i]), len(idx.in[i]))

				for j := 0; j < len(tt.want.in[i]); j++ {
					assert.Equal(t, tt.want.in[i][j].Interface(), idx.in[i][j].Interface())
				}
			}
			assert.Equal(t, len(tt.want.out), len(idx.out))
			for i := 0; i < len(tt.want.out); i++ {
				assert.Equal(t, len(tt.want.out[i]), len(idx.out[i]))

				for j := 0; j < len(tt.want.out[i]); j++ {
					assert.Equal(t, tt.want.out[i][j].Interface(), idx.out[i][j].Interface())
				}
			}
		})
	}
}

func Test_callsIndex_MockedOutFor(t *testing.T) {
	type fields struct {
		in  [][]reflect.Value
		out [][]reflect.Value
	}
	type args struct {
		in []reflect.Value
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []reflect.Value
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &callsIndex{
				in:  tt.fields.in,
				out: tt.fields.out,
			}
			got, err := i.MockedOutFor(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("callsIndex.MockedOutFor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("callsIndex.MockedOutFor() = %v, want %v", got, tt.want)
			}
		})
	}
}
