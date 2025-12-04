package repository

import (
	"context"

	"github.com/gdugdh24/mpit2026-backend/internal/domain"
)

type MessageRepository interface {
	Create(ctx context.Context, message *domain.Message) error
	GetByID(ctx context.Context, id int) (*domain.Message, error)
	GetMatchMessages(ctx context.Context, matchID int, limit, offset int) ([]*domain.Message, error)
	MarkAsRead(ctx context.Context, messageID int) error
	GetUnreadCount(ctx context.Context, userID int) (int, error)
	Delete(ctx context.Context, id int) error
}
