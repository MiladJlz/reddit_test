package api

import (
	"fmt"
	"net/http"
)

type Error struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

func (e Error) Error() string {
	return e.Err
}

func NewError(code int, err string) Error {
	return Error{
		Code: code,
		Err:  err,
	}
}

func ErrResourceNotFound(res string) Error {
	return Error{
		Code: http.StatusNotFound,
		Err:  res + " resource not found",
	}
}

func ErrBadRequest() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err:  "invalid JSON request",
	}
}

func ErrInvalidID() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err:  "invalid id given",
	}
}

func ErrInvalidPathParam() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err:  "id must be a UUID",
	}
}
func ErrFailedCreatingUUID(res string) Error {
	return Error{
		Code: http.StatusInternalServerError,
		Err:  fmt.Sprintf("failed creating uuid for %s", res),
	}
}

func ErrFailedInsertingData(res string) Error {
	return Error{
		Code: http.StatusInternalServerError,
		Err:  fmt.Sprintf("failed inserting %s to DB", res),
	}
}
func ErrFailedGettingData(res string) Error {
	return Error{
		Code: http.StatusInternalServerError,
		Err:  fmt.Sprintf("failed Getting %s from DB", res),
	}
}
func ErrVoteConflict() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err:  "user already voted",
	}
}
func ErrInvalidRequestBody() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err:  "invalid request body",
	}
}

func ErrInvalidQueryParam() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err:  "invalid query parameter",
	}
}
