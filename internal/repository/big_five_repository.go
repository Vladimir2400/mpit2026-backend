package repository

import (
	"context"

	"github.com/gdugdh24/mpit2026-backend/internal/domain"
)

type BigFiveRepository interface {
	Create(ctx context.Context, result *domain.BigFiveResult) error
	GetByID(ctx context.Context, id int) (*domain.BigFiveResult, error)
	GetByUserID(ctx context.Context, userID int) (*domain.BigFiveResult, error)
	Update(ctx context.Context, result *domain.BigFiveResult) error
	Delete(ctx context.Context, id int) error
}
