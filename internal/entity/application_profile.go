package entity

import (
	"time"

	"github.com/google/uuid"
)

type ApplicationProfile struct {
	ID            uuid.UUID
	ApplicationID uuid.UUID
	Version       int64
	CreatedAt     time.Time
	GraphID       string
}
