package models

import (
	"github.com/google/uuid"
	"time"
)

type Lesson struct {
	ID        uuid.UUID
	TutorId   uuid.UUID
	StudentId uuid.UUID
	Price     int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
