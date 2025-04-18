package entity

import (
	"time"

	"github.com/google/uuid"
)

type Application struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
}
