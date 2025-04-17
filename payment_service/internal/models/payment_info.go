package models

import "github.com/google/uuid"

type PaymentInfo struct {
	LessonID       uuid.UUID
	PriceRUB       uint64
	PaymentDetails string
}
