package models

import (
	"github.com/google/uuid"
)

type UserContext struct {
	UserID   uuid.UUID
	UserRole Role
}
