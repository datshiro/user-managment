package helper

import (
	"reflect"
	"testing"
  "net/http"
  "errors"
)

func TestAPIError_WithDetails(t *testing.T) {
	type fields struct {
		StatusCode int
		Message    any
		details    error
	}
	type args struct {
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   APIError
	}{
    {
      "test successful",
      fields{
        StatusCode: http.StatusBadRequest,
        Message: "invalid JSON request data",
        details: nil,
      },
      args{
        err: errors.New("this is details error"),
      },
      APIError{
        StatusCode: http.StatusBadRequest,
        Message: "invalid JSON request data",
        details: errors.New("this is details error"),
      },
    },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := APIError{
				StatusCode: tt.fields.StatusCode,
				Message:    tt.fields.Message,
				details:    tt.fields.details,
			}
			if got := e.WithDetails(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("APIError.WithDetails() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIError_Error(t *testing.T) {
	type fields struct {
		StatusCode int
		Message    any
		details    error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
    {
      "Print error",
      fields{
        StatusCode: http.StatusInternalServerError,
        Message: "this is error message",
        details: nil,
      },
      "api error: 500; this is error message",
    },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := APIError{
				StatusCode: tt.fields.StatusCode,
				Message:    tt.fields.Message,
				details:    tt.fields.details,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("APIError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAPIError(t *testing.T) {
	type args struct {
		statusCode int
		err        error
	}
	tests := []struct {
		name string
		args args
		want APIError
	}{
    {
      "create new error",
      args{
        statusCode: http.StatusBadRequest,
        err: errors.New("invalid JSON request data"),
      },
      APIError{
        StatusCode: http.StatusBadRequest,
        Message: "invalid JSON request data",
        details: nil,
      },
    },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAPIError(tt.args.statusCode, tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAPIError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInvalidJSON(t *testing.T) {
	tests := []struct {
		name string
		want APIError
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InvalidJSON(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InvalidJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInvalidValidationErrors(t *testing.T) {
	type args struct {
		errs map[string]error
	}
	tests := []struct {
		name string
		args args
		want APIError
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InvalidValidationErrors(tt.args.errs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InvalidValidationErrors() = %v, want %v", got, tt.want)
			}
		})
	}
}
