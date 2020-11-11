package apperr

import "fmt"

type Error struct {
	Message       string `json:"message"`
	Hint          string `json:"hint"`
	InternalError error  `json:"-"`
}

func (e Error) Error() string {
	return fmt.Sprintf("app_eror: %s. internal_err: %v", e.Message, e.InternalError)
}

func New(message, hint string) error {
	return Error{
		Message: message,
		Hint:    hint,
	}
}

func With(err error, message string, hint string) error {
	// apiErr, ok := err.(Error)
	// if ok {
	// 	return apiErr
	// }
	return Error{
		Message:       message,
		Hint:          hint,
		InternalError: err,
	}
}
