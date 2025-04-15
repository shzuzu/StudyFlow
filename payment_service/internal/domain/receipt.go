package domain

import (
	"time"

	"github.com/google/uuid"
)

type PaymentReceipt struct {
	ID         uuid.UUID
	LessonID   uuid.UUID
	FileID     uuid.UUID
	IsVerified bool
	CreatedAt  time.Time
	EditedAt   *time.Time
}
