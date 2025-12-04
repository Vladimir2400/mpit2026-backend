package repository

import (
	"context"

	"github.com/gdugdh24/mpit2026-backend/internal/domain"
)

type SwipeRepository interface {
	Create(ctx context.Context, swipe *domain.Swipe) error
	GetByID(ctx context.Context, id int) (*domain.Swipe, error)
	GetByUsers(ctx context.Context, swiperID, swipedID int) (*domain.Swipe, error)
	GetUserSwipes(ctx context.Context, userID int, limit, offset int) ([]*domain.Swipe, error)
	GetLikesReceived(ctx context.Context, userID int, limit, offset int) ([]*domain.Swipe, error)
	CheckMutualLike(ctx context.Context, user1ID, user2ID int) (bool, error)
}
