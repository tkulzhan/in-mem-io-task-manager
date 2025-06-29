package errors

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
		Code:    404,
	}
}

func NewInternalError(message string) *Error {
	return &Error{
		Message: message,
		Code:    500,
	}
}

func NewBadRequestError(message string) *Error {
	return &Error{
		Message: message,
		Code:    400,
	}
}

func (e *Error) Error() string {
	return e.Message
}