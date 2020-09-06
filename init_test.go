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
			name: "init with error list",
			args: args{
				params: []interface{}{[]Error{{Message: "another error"}}},
			},
			want: Error{
				Errors: []*Error{
					{Message: "another error"},
				},
			},
		},
		{
			name: "init with error pointer list",
			args: args{
				params: []interface{}{[]*Error{{Message: "another error"}}},
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

func TestEt(t *testing.T) {
	type args struct {
		target string
		args   []interface{}
	}
	tests := []struct {
		name string
		args args
		want Error
	}{
		{
			name: "target is not empty",
			args: args{
				target: "new-target",
				args:   nil,
			},
			want: Error{
				Target: "new-target",
			},
		},
		{
			name: "target is empty",
			args: args{
				target: "",
				args:   nil,
			},
			want: Error{
				Target: "",
			},
		},
		{
			name: "target with another target",
			args: args{
				target: "target",
				args: []interface{}{
					Target("new-target"),
				},
			},
			want: Error{
				Target: "target",
			},
		},
		{
			name: "target with message",
			args: args{
				target: "target",
				args: []interface{}{
					"message",
				},
			},
			want: Error{
				Target:  "target",
				Message: "message",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Et(tt.args.target, tt.args.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Et() = %v, want %v", got, tt.want)
			}
		})
	}
}
