package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gdugdh24/mpit2026-backend/internal/domain"
	"github.com/gdugdh24/mpit2026-backend/internal/repository"
	"github.com/jmoiron/sqlx"
)

type messageRepository struct {
	db *sqlx.DB
}

func NewMessageRepository(db *sqlx.DB) repository.MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(ctx context.Context, message *domain.Message) error {
	query := `
		INSERT INTO messages (match_id, sender_id, content, is_read)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	return r.db.QueryRowContext(
		ctx, query,
		message.MatchID, message.SenderID, message.Content, message.IsRead,
	).Scan(&message.ID, &message.CreatedAt)
}

func (r *messageRepository) GetByID(ctx context.Context, id int) (*domain.Message, error) {
	var message domain.Message
	query := `SELECT * FROM messages WHERE id = $1`
	err := r.db.GetContext(ctx, &message, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrMessageNotFound
		}
		return nil, err
	}
	return &message, nil
}

func (r *messageRepository) GetMatchMessages(ctx context.Context, matchID int, limit, offset int) ([]*domain.Message, error) {
	var messages []*domain.Message
	query := `
		SELECT * FROM messages
		WHERE match_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	err := r.db.SelectContext(ctx, &messages, query, matchID, limit, offset)
	return messages, err
}

func (r *messageRepository) MarkAsRead(ctx context.Context, messageID int) error {
	query := `UPDATE messages SET is_read = true WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, messageID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrMessageNotFound
	}
	return nil
}

func (r *messageRepository) GetUnreadCount(ctx context.Context, userID int) (int, error) {
	var count int
	query := `
		SELECT COUNT(*)
		FROM messages m
		JOIN matches ma ON m.match_id = ma.id
		WHERE (ma.user1_id = $1 OR ma.user2_id = $1)
		AND m.sender_id != $1
		AND m.is_read = false
	`
	err := r.db.GetContext(ctx, &count, query, userID)
	return count, err
}

func (r *messageRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM messages WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrMessageNotFound
	}
	return nil
}
