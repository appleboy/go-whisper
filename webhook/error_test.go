package webhook

import (
	"errors"
	"testing"
)

func TestRequestError_Error(t *testing.T) {
	type fields struct {
		HTTPStatusCode int
		Err            error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test error message",
			fields: fields{
				HTTPStatusCode: 404,
				Err:            errors.New("not found"),
			},
			want: "error, status code: 404, message: not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &RequestError{
				HTTPStatusCode: tt.fields.HTTPStatusCode,
				Err:            tt.fields.Err,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("RequestError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
