package errors1

import "errors"

var (
	ErrLessonNotFound    = errors.New("lesson not found")
	ErrPermissionDenied  = errors.New("permission was denied")
	ErrInvalidPayment    = errors.New("invalid payment")
	ErrReceiptNotFound   = errors.New("receipt was not found")
	ErrLessonAlreadyPaid = errors.New("lesson was already paid")
	ErrInvalidArgument   = errors.New("invalid argument")
	ErrAlreadyExists     = errors.New("receipt already exists")
)
