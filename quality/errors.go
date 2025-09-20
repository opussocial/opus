package quality

import "fmt"

type Error struct {
	msg string
}

func (e *Error) Error() string {
	return e.msg
}

func (e *Error) WithDetail(detail string) error {
	return fmt.Errorf("%w: %s", e, detail)
}

func (e *Error) Unwrap() error {
	return nil
}

func NewError(msg string) *Error {
	return &Error{msg: msg}
}
