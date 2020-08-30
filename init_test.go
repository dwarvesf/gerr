package gerr

import (
	"reflect"
	"testing"
)

func TestE(t *testing.T) {
	type args struct {
		params []interface{}
	}
	tests := []struct {
		name string
		args args
		want Error
	}{
		{
			name: "init with target",
			args: args{
				params: []interface{}{Target("target-1")},
			},
			want: Error{
				Target: "target-1",
			},
		},
		{
			name: "init with message",
			args: args{
				params: []interface{}{Message("err message")},
			},
			want: Error{
				Message: "err message",
			},
		},
		{
			name: "init with code",
			args: args{
				params: []interface{}{Code(400)},
			},
			want: Error{
				Code: 400,
			},
		},
		{
			name: "init with error",
			args: args{
				params: []interface{}{Error{Message: "another error"}},
			},
			want: Error{
				Errors: []*Error{
					{Message: "another error"},
				},
			},
		},
		{
			name: "init with error pointer",
			args: args{
				params: []interface{}{&Error{Message: "another error"}},
			},
			want: Error{
				Errors: []*Error{
					{Message: "another error"},
				},
			},
		},
		{
			name: "init with unknown field",
			args: args{
				params: []interface{}{float64(3)},
			},
			want: Error{
				Errors: []*Error{
					{Message: "unknown type float64, value 3 in error call"},
				},
			},
		},
		{
			name: "init with name, code",
			args: args{
				params: []interface{}{"str", 400},
			},
			want: Error{
				Message: "str", Code: 400,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := E(tt.args.params...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("E() = %v, want %v", got, tt.want)
			}
		})
	}
}
