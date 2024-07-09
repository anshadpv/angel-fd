package errors

import (
	"bytes"
	"errors"
	"fmt"
)

const (
	ErrInternal = "internal"
	Prefix      = "FD-CORE-"
)

type Error struct {
	Code     string
	Message  string
	HTTPCode int
	Op       string
	Err      error
	Log      bool
}

type ErrorBuilder interface {
	Code(code string) ErrorBuilder
	Msg(message string) ErrorBuilder
	Msgf(format string, v ...interface{}) ErrorBuilder
	HTTPCode(httpCode int) ErrorBuilder
	Op(op string) ErrorBuilder
	Wrap(err error) ErrorBuilder
	Build() *Error
}

type errorBuilder struct {
	code     string
	message  string
	httpCode int
	op       string
	err      error
	log      bool
}

func (builder *errorBuilder) Code(code string) ErrorBuilder {
	builder.code = Prefix + code
	return builder
}

func (builder *errorBuilder) Msg(message string) ErrorBuilder {
	builder.message = message
	return builder
}

func (builder *errorBuilder) Msgf(format string, v ...interface{}) ErrorBuilder {
	if builder == nil {
		return builder
	}
	return builder.Msg(fmt.Sprintf(format, v...))
}

func (builder *errorBuilder) HTTPCode(httpCode int) ErrorBuilder {
	builder.httpCode = httpCode
	return builder
}

func (builder *errorBuilder) Op(op string) ErrorBuilder {
	builder.op = op
	return builder
}

func (builder *errorBuilder) Log(log bool) ErrorBuilder {
	builder.log = log
	return builder
}

func (builder *errorBuilder) Wrap(err error) ErrorBuilder {
	builder.err = err
	return builder
}

func (builder *errorBuilder) Build() *Error {
	return &Error{
		Code:     builder.code,
		Message:  builder.message,
		HTTPCode: builder.httpCode,
		Op:       builder.op,
		Err:      builder.err,
		Log:      builder.log,
	}
}

func New() ErrorBuilder {
	return &errorBuilder{log: true}
}

func NewWithHTTP(internalCode string, httpCode int, message string) *Error {
	err := &errorBuilder{code: internalCode, httpCode: httpCode, message: message}
	return err.Build()
}

func UseWithHTTP(internalCode string, httpCode int, message string, log bool) *errorBuilder {
	return &errorBuilder{code: Prefix + internalCode, httpCode: httpCode, message: message, log: log}
}

func NewWith(internalCode string, message string) *Error {
	err := &errorBuilder{code: internalCode, message: message}
	return err.Build()
}

func From(err error) ErrorBuilder {
	var builder errorBuilder
	var e *Error
	if errors.As(err, &e) {
		builder.code = e.Code
		builder.message = e.Message
		builder.httpCode = e.HTTPCode
		builder.op = e.Op
		builder.err = e.Err
	} else {
		builder.err = err
	}
	return &builder
}

func (e *Error) Error() string {
	var buf bytes.Buffer

	// Print the current operation in our stack, if any
	if e.Op != "" {
		fmt.Fprintf(&buf, "%s:", e.Op)
	}

	// If wrapping an error, print its Error() message.
	// Otherwise print the error code and message
	if e.Err != nil {
		buf.WriteString(e.Err.Error())
	} else {
		if e.Message != "" {
			buf.WriteString(e.Message)
		}
	}

	return buf.String()
}

func ErrorCode(err error) string {
	if err == nil {
		return ""
	}
	var e *Error
	if ok := errors.As(err, &e); ok && e.Code != "" {
		return e.Code
	} else if ok && e.Err != nil {
		return ErrorCode(e.Err)
	}

	return ErrInternal
}

func ErrorMessage(err error) string {
	if err == nil {
		return ""
	}
	var e *Error
	if ok := errors.As(err, &e); ok && e.Message != "" {
		return e.Message
	} else if ok && e.Err != nil {
		return ErrorMessage(e.Err)
	}

	return "An internal error has occurred. Please contact support"
}
