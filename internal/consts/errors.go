package consts

import (
	"fmt"
	"net/http"
)

// Custom error type for internal service
// this struct must be initialize via 
type CustomError struct {
	root     error
	err      error
	args     []interface{}
	tags     map[string]interface{}
	httpCode int
}

func NewError(err error, args ...interface{}) CustomError {
	return CustomError{root: err, err: err, tags: map[string]interface{}{}}
}

// Error message is combination of message, args and tags
func (e CustomError) Error() string {
	var msg string
	re, ok := e.root.(CustomError)
	if ok {
		msg += fmt.Sprintf("%v; ", re.Error())
	}
	if e.err != nil {
		msg += fmt.Sprintf("%v; ", e.err.Error())
	}
	for _, arg := range e.args {
		msg += fmt.Sprintf(" %+v; ", arg)
	}

	return msg
}

// Add tag to error for additional info, which useful for log and tracing
func (e CustomError) WithTag(key string, value interface{}) CustomError {
	e.tags[key] = value
	return e
}

// set root cause error
func (e CustomError) WithRootCause(err error) CustomError {
	e.root = err
	return e
}

func (e CustomError) Details() string {
	var msg string
	if e.root != nil {
		msg += e.root.Error()
	}
	for tag, value := range e.tags {
		msg += fmt.Sprintf("; %s=%+v", tag, value)
	}
	return msg
}

// Set HTTP Status Code which will set to header in response
func (e CustomError) WithHttpCode(code int) CustomError {
	e.httpCode = code
	return e
}

// Set HTTP Status Code which will set to header in response
func (e CustomError) GetCode() int {
	if e.httpCode == 0 {
		return http.StatusBadRequest // 400 as default
	}
	return e.httpCode
}

// Init error with error type
func withError(err error, args ...interface{}) CustomError {
	return NewError(err, args...)
}

// Init error with string message
func withMessage(msg string, args ...interface{}) CustomError {
	return NewError(fmt.Errorf(msg), args...)
}

var (
	ErrInvalidRequest        = withMessage("Invalid request")
	ErrCreateFailure         = withMessage("Create failure")
	ErrLoginFailure          = withMessage("Login Failure")
	ErrDataNotFound          = withMessage("Requesting resource not found")
	ErrMissingCredentialInfo = withMessage("account info cannot be empty, please provide username, email or phone number")
	ErrPasswordCannotBeEmpty = withMessage("password cannot be empty")
)
