package helper

import (
	"fmt"
	"net/http"
)

type APIError struct {
	StatusCode int   `json:"status_code"`
	Message    any   `json:"message"`
	details    error 
}

func (e APIError) WithDetails(err error) APIError {
	e.details = err
	return e
}

func (e APIError) Error() string {
	return fmt.Sprintf("api error: %d %s ", e.StatusCode, e.Message)
}

func NewAPIError(statusCode int, err error) APIError {
	return APIError{StatusCode: statusCode, Message: err.Error()}
}

func InvalidJSON() APIError {
	return NewAPIError(http.StatusBadRequest, fmt.Errorf("invalid request JSON data"))
}

func InvalidValidationErrors(errs map[string]error) APIError {
	return APIError{StatusCode: http.StatusUnprocessableEntity, Message: errs}
}
