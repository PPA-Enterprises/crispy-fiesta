package common

import (
	"net/http"
)

type ResponseError struct {
	Code int
	reason string
}

func (e *ResponseError) Error() string {
	return e.reason
}

func EmailAlreadyExistsError() *ResponseError {
	return &ResponseError{
		Code: http.StatusConflict,
		reason: "Email Already exists",
	}
}

func DatabaseError(err error) *ResponseError {
	return &ResponseError{
		Code: http.StatusInternalServerError,
		reason: err.Error(),
	}
}

func UidTypeAssertionError() *ResponseError {
	return &ResponseError{
		Code: http.StatusInternalServerError,
		reason: "ObjectId type assertion failed",
	}
}

func ArgonHashError(err error) *ResponseError {
	return &ResponseError{
		Code: http.StatusInternalServerError,
		reason: err.Error(),
	}
}
