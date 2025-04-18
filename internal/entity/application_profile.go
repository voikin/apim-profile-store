package entity

import (
	"time"

	"github.com/google/uuid"
)

type ApplicationProfile struct {
	ID            uuid.UUID
	ApplicationID uuid.UUID
	Version       uint32
	GraphID       uuid.UUID
	CreatedAt     time.Time
}
