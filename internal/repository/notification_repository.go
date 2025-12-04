package repository

import (
	"context"

	"github.com/gdugdh24/mpit2026-backend/internal/domain"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification *domain.Notification) error
	GetByID(ctx context.Context, id int) (*domain.Notification, error)
	GetUserNotifications(ctx context.Context, userID int, limit, offset int) ([]*domain.Notification, error)
	MarkAsRead(ctx context.Context, notificationID int) error
	MarkAllAsRead(ctx context.Context, userID int) error
	GetUnreadCount(ctx context.Context, userID int) (int, error)
	Delete(ctx context.Context, id int) error
}
