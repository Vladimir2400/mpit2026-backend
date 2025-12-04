package repository

import (
	"context"

	"github.com/gdugdh24/mpit2026-backend/internal/domain"
)

type SessionRepository interface {
	Create(ctx context.Context, session *domain.Session) error
	GetByToken(ctx context.Context, token string) (*domain.Session, error)
	GetByUserID(ctx context.Context, userID int) ([]*domain.Session, error)
	Delete(ctx context.Context, id int) error
	DeleteByToken(ctx context.Context, token string) error
	DeleteExpired(ctx context.Context) error
	DeleteByUserID(ctx context.Context, userID int) error
}
