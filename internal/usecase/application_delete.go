package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (u *Usecase) DeleteApplication(ctx context.Context, id uuid.UUID) error {
	err := u.postgresRepo.DeleteApplication(ctx, id)
	if err != nil {
		return fmt.Errorf("postgresRepo.DeleteApplication: %w", err)
	}

	return nil
}
