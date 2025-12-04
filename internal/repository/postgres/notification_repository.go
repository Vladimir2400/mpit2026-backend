package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gdugdh24/mpit2026-backend/internal/domain"
	"github.com/gdugdh24/mpit2026-backend/internal/repository"
	"github.com/jmoiron/sqlx"
)

type notificationRepository struct {
	db *sqlx.DB
}

func NewNotificationRepository(db *sqlx.DB) repository.NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(ctx context.Context, notification *domain.Notification) error {
	query := `
		INSERT INTO notifications (user_id, content, is_read)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	return r.db.QueryRowContext(
		ctx, query,
		notification.UserID, notification.Content, notification.IsRead,
	).Scan(&notification.ID, &notification.CreatedAt)
}

func (r *notificationRepository) GetByID(ctx context.Context, id int) (*domain.Notification, error) {
	var notification domain.Notification
	query := `SELECT * FROM notifications WHERE id = $1`
	err := r.db.GetContext(ctx, &notification, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("notification not found")
		}
		return nil, err
	}
	return &notification, nil
}

func (r *notificationRepository) GetUserNotifications(ctx context.Context, userID int, limit, offset int) ([]*domain.Notification, error) {
	var notifications []*domain.Notification
	query := `
		SELECT * FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	err := r.db.SelectContext(ctx, &notifications, query, userID, limit, offset)
	return notifications, err
}

func (r *notificationRepository) MarkAsRead(ctx context.Context, notificationID int) error {
	query := `UPDATE notifications SET is_read = true WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, notificationID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("notification not found")
	}
	return nil
}

func (r *notificationRepository) MarkAllAsRead(ctx context.Context, userID int) error {
	query := `UPDATE notifications SET is_read = true WHERE user_id = $1 AND is_read = false`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *notificationRepository) GetUnreadCount(ctx context.Context, userID int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = false`
	err := r.db.GetContext(ctx, &count, query, userID)
	return count, err
}

func (r *notificationRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM notifications WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("notification not found")
	}
	return nil
}
