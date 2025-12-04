package repository

import (
	"context"

	"github.com/gdugdh24/mpit2026-backend/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id int) (*domain.User, error)
	GetByVKID(ctx context.Context, vkID int) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id int) error
	UpdateOnlineStatus(ctx context.Context, userID int, isOnline bool) error
	GetOnlineUsers(ctx context.Context, limit, offset int) ([]*domain.User, error)
}
