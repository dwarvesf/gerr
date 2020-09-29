package gerr

import (
	"reflect"
	"testing"
)

func Test_doMakeErrResponse(t *testing.T) {
	type args struct {
		err Error
	}
	tests := []struct {
		name string
		args args
		want ErrResponse
	}{
		{
			name: "will return simple error",
			args: args{err: Error{
				Message: "message error",
				Target:  "",
				Errors: []error{
					&Error{
						Target:  "field1",
						Message: "error field1",
					},
					&Error{
						Target:  "field2",
						Message: "error field2",
					},
				},
			}},
			want: ErrResponse{
				Message: "message error",
				Errors: map[string]interface{}{
					"field1": []interface{}{"error field1"},
					"field2": []interface{}{"error field2"},
				},
			},
		},
		{
			name: "will return error list",
			args: args{
				err: Error{
					Message: "message error",
					Target:  "",
					Errors: []error{
						&Error{
							Target:  "items",
							Message: "items got error",
							Errors: []error{
								&Error{
									Target: "0",
									Errors: []error{
										&Error{
											Target:  "amount",
											Message: "out of stock",
										},
									},
								},
								&Error{
									Target: "1",
									Errors: []error{
										&Error{
											Target:  "id",
											Message: "invalid",
											Errors: []error{
												&Error{
													Message: "not found",
												},
												&Error{
													Message: "invalid",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: ErrResponse{
				Message: "message error",
				Errors: map[string]interface{}{
					"items": map[string]interface{}{
						"0": map[string]interface{}{
							"amount": []interface{}{"out of stock"},
						},
						"1": map[string]interface{}{
							"id": []interface{}{"not found", "invalid"},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := doMakeErrResponse(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("doMakeErrResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
