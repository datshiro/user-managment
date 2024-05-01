package consts

import (
	"fmt"
)

// Custom error type for internal service
// this struct must be initialize via NewCakeError
type CakeError struct {
	root error
	err  error
	args []interface{}
	tags map[string]interface{}
}

func NewCakeError(err error, args ...interface{}) CakeError {
	return CakeError{root: err, tags: map[string]interface{}{}}
}

// Error message is combination of message, args and tags
func (e CakeError) Error() string {
	msg := e.err.Error()
	for _, arg := range e.args {
		msg += fmt.Sprintf("; %+v", arg)
	}

	for tag, value := range e.tags {
		msg += fmt.Sprintf("; %s=%+v", tag, value)
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

// Init error with error type
func withError(err error, args ...interface{}) CakeError {
	return NewCakeError(err, args...)
}

// Init error with string message
func withMessage(msg string, args ...interface{}) CakeError {
	return NewCakeError(fmt.Errorf(msg), args...)
}

func (e CakeError) Details() error {
  return e.root
}

var (
	ErrInvalidRequest = withMessage("Invalid Request")
)
