package errdefs

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrPermissionDenied     = errors.New("permission was denied")
	ErrInvalidArgument      = errors.New("invalid argument")
	ErrNotFound             = errors.New("not found")
	ErrInvalidPayment       = errors.New("invalid payment")
	ErrAlreadyExists        = errors.New("already exists")
	InternalError           = errors.New("internal error")
	ErrLessonNotFound       = errors.New("lesson not found")
	ErrLessonAlreadyPaid    = errors.New("lesson already paid")
	ErrReceiptAlreadyExists = errors.New("receipt already exists")
)

func MapError(err error) error {
	if err == nil {
		return nil
	}

	if s, ok := status.FromError(err); ok {
		switch s.Code() {
		case codes.NotFound:
			return ErrNotFound
		case codes.InvalidArgument:
			return ErrInvalidArgument
		case codes.AlreadyExists:
			return ErrAlreadyExists
		case codes.PermissionDenied:
			return ErrPermissionDenied
		case codes.FailedPrecondition:
			return ErrInvalidPayment
		case codes.Unavailable, codes.DeadlineExceeded, codes.Aborted:
			return ErrNotFound
		default:
			return InternalError
		}
	}

	switch {
	case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
		return ErrNotFound
	case errors.Is(err, ErrNotFound):
		return ErrNotFound
	case errors.Is(err, ErrInvalidArgument):
		return ErrInvalidArgument
	case errors.Is(err, ErrAlreadyExists):
		return ErrAlreadyExists
	case errors.Is(err, ErrPermissionDenied):
		return ErrPermissionDenied
	case errors.Is(err, ErrInvalidPayment):
		return ErrInvalidPayment
	case errors.Is(err, ErrLessonNotFound):
		return ErrLessonNotFound
	case errors.Is(err, ErrLessonAlreadyPaid):
		return ErrLessonAlreadyPaid
	case errors.Is(err, ErrReceiptAlreadyExists):
		return ErrReceiptAlreadyExists
	default:
		return InternalError
	}
}
