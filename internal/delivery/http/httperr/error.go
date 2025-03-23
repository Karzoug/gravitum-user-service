package httperr

import "net/http"

// Error is a http transport level error.
type Error struct {
	msg  string
	code int
}

func NewError(msg string, code int) Error {
	return Error{
		msg:  msg,
		code: code,
	}
}

func NewInternalError(err error) Error {
	return Error{
		msg:  http.StatusText(http.StatusInternalServerError),
		code: http.StatusInternalServerError,
	}
}

func (e Error) HTTPStatus() (int, string) {
	return e.code, e.msg
}

func (e Error) Error() string {
	return e.msg
}
