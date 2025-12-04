package repository

import (
	"context"

	"github.com/gdugdh24/mpit2026-backend/internal/domain"
)

type ProfileRepository interface {
	Create(ctx context.Context, profile *domain.Profile) error
	GetByID(ctx context.Context, id int) (*domain.Profile, error)
	GetByUserID(ctx context.Context, userID int) (*domain.Profile, error)
	Update(ctx context.Context, profile *domain.Profile) error
	Delete(ctx context.Context, id int) error
	UpdateOnboardingStatus(ctx context.Context, userID int, isComplete bool) error
	SearchProfiles(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*domain.Profile, error)
}
