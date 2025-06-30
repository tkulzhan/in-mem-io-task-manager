package errors

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func New(message string, code int) *Error {
	return &Error{
		Message: message,
		Code:    code,
	}
}

func NewNotFoundError(message string) *Error {
	return &Error{
		Message: message,
		Code:    http.StatusNotFound,
	}
}

func NewInternalError(message string) *Error {
	return &Error{
		Message: message,
		Code:    http.StatusInternalServerError,
	}
}

func NewBadRequestError(message string) *Error {
	return &Error{
		Message: message,
		Code:    http.StatusBadRequest,
	}
}

func (e Error) Error() string {
	return e.Message
}

func HandleError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	if err == nil {
		cerr := NewInternalError("nil error passed to HandleError")
		writeJSONError(w, cerr, http.StatusInternalServerError)
		return
	}

	var cerr *Error
	if errors.As(err, &cerr) && cerr != nil {
		writeJSONError(w, *cerr, cerr.Code)
		return
	}

	unexpectedErr := NewInternalError(err.Error())
	writeJSONError(w, unexpectedErr, http.StatusInternalServerError)
}

func writeJSONError(w http.ResponseWriter, err error, status int) {
	w.WriteHeader(status)
	if encodeErr := json.NewEncoder(w).Encode(err); encodeErr != nil {
		http.Error(w, encodeErr.Error(), http.StatusInternalServerError)
	}
}
