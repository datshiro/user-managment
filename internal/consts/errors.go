package consts

import (
	"fmt"
	"net/http"
)

// Custom error type for internal service
// this struct must be initialize via NewCakeError
type CakeError struct {
	root     error
	err      error
	args     []interface{}
	tags     map[string]interface{}
	httpCode int
}

func NewCakeError(err error, args ...interface{}) CakeError {
  return CakeError{root: err, err: err, tags: map[string]interface{}{}}
}

// Error message is combination of message, args and tags
func (e CakeError) Error() string {
	var msg string
	if e.err != nil {
		msg = e.err.Error()
	}
	for _, arg := range e.args {
		msg += fmt.Sprintf("; %+v", arg)
	}

	return msg
}

// Add tag to error for additional info, which useful for log and tracing
func (e CakeError) WithTag(key string, value interface{}) CakeError {
	e.tags[key] = value
	return e
}

// set root cause error
func (e CakeError) WithRootCause(err error) CakeError {
	e.root = err
	return e
}

func (e CakeError) Details() string {
	var msg string 
  if e.root != nil {
    msg  = e.root.Error()
  }
	for tag, value := range e.tags {
		msg += fmt.Sprintf("; %s=%+v", tag, value)
	}
	return msg
}

// Set HTTP Status Code which will set to header in response
func (e CakeError) WithHttpCode(code int) CakeError {
  e.httpCode = code
	return e
}

// Set HTTP Status Code which will set to header in response
func (e CakeError) GetCode() int {
  if e.httpCode == 0{
    return http.StatusBadRequest // 400 as default
  }
	return e.httpCode
}

// Init error with error type
func withError(err error, args ...interface{}) CakeError {
	return NewCakeError(err, args...)
}

// Init error with string message
func withMessage(msg string, args ...interface{}) CakeError {
	return NewCakeError(fmt.Errorf(msg), args...)
}


var (
	ErrInvalidRequest = withMessage("Invalid Request")
	ErrCreateFailure  = withMessage("Create failure")
	ErrLoginFailure   = withMessage("Provided login credential doesn't support")
	ErrDataNotFound   = withMessage("Requesting resource not found")
)
