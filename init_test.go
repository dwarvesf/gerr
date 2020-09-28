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
				params: []interface{}{
					Target("target-1"),
					Op("TestE.func1"),
				},
			},
			want: Error{
				Target: "target-1",
				Op:     "TestE.func1",
			},
		},
		{
			name: "init with message",
			args: args{
				params: []interface{}{
					Message("err message"),
					Op("TestE.func1"),
				},
			},
			want: Error{
				Message: "err message",
				Op:      "TestE.func1",
			},
		},
		{
			name: "init with code",
			args: args{
				params: []interface{}{
					Op("TestE.func1"),
					Code(400),
				},
			},
			want: Error{
				Code:    400,
				Message: "Bad Request",
				Op:      "TestE.func1",
			},
		},
		{
			name: "init with error",
			args: args{
				params: []interface{}{
					"error",
					400,
					Error{Message: "another error"},
					Op("TestE.func1"),
				},
			},
			want: Error{
				Message: "error",
				Code:    400,
				Op:      "TestE.func1",
				Errors: []*Error{
					{Message: "another error"},
				},
			},
		},
		{
			name: "init with error pointer",
			args: args{
				params: []interface{}{
					"error",
					400,
					&Error{Message: "another error"},
					Op("TestE.func1"),
				},
			},
			want: Error{
				Message: "error",
				Code:    400,
				Op:      "TestE.func1",
				Errors: []*Error{
					{Message: "another error"},
				},
			},
		},
		{
			name: "init with error list",
			args: args{
				params: []interface{}{Op("TestE.func1"), []Error{{Message: "another error"}}},
			},
			want: Error{
				Op: "TestE.func1",
				Errors: []*Error{
					{Message: "another error"},
				},
			},
		},
		{
			name: "init with error pointer list",
			args: args{
				params: []interface{}{Op("TestE.func1"), []*Error{{Message: "another error"}}},
			},
			want: Error{
				Op: "TestE.func1",
				Errors: []*Error{
					{Message: "another error"},
				},
			},
		},
		{
			name: "init with unknown field",
			args: args{
				params: []interface{}{Op("TestE.func1"), float64(3)},
			},
			want: Error{
				Op: "TestE.func1",
				Errors: []*Error{
					{Message: "unknown type float64, value 3 in error call"},
				},
			},
		},
		{
			name: "init with name, code",
			args: args{
				params: []interface{}{Op("TestE.func1"), "str", 400},
			},
			want: Error{
				Op:      "TestE.func1",
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
				args: []interface{}{
					Op("TestEt.func1"),
				},
			},
			want: Error{
				Op:     "TestEt.func1",
				Target: "new-target",
			},
		},
		{
			name: "target is empty",
			args: args{
				target: "",
				args: []interface{}{
					Op("TestEt.func1"),
				},
			},
			want: Error{
				Op:     "TestEt.func1",
				Target: "",
			},
		},
		{
			name: "target with another target",
			args: args{
				target: "target",
				args: []interface{}{
					Op("TestEt.func1"),
					Target("new-target"),
				},
			},
			want: Error{
				Op:     "TestEt.func1",
				Target: "target",
			},
		},
		{
			name: "target with message",
			args: args{
				target: "target",
				args: []interface{}{
					Op("TestEt.func1"),
					"message",
				},
			},
			want: Error{
				Op:      "TestEt.func1",
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
