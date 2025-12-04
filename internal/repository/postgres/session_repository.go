package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gdugdh24/mpit2026-backend/internal/domain"
	"github.com/gdugdh24/mpit2026-backend/internal/repository"
	"github.com/jmoiron/sqlx"
)

type sessionRepository struct {
	db *sqlx.DB
}

func NewSessionRepository(db *sqlx.DB) repository.SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) Create(ctx context.Context, session *domain.Session) error {
	query := `
		INSERT INTO sessions (user_id, token, device_info, ip_address, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	return r.db.QueryRowContext(
		ctx, query,
		session.UserID, session.Token, session.DeviceInfo,
		session.IPAddress, session.ExpiresAt,
	).Scan(&session.ID, &session.CreatedAt)
}

func (r *sessionRepository) GetByToken(ctx context.Context, token string) (*domain.Session, error) {
	var session domain.Session
	query := `SELECT * FROM sessions WHERE token = $1`
	err := r.db.GetContext(ctx, &session, query, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrSessionNotFound
		}
		return nil, err
	}
	if session.IsExpired() {
		return nil, domain.ErrSessionExpired
	}
	return &session, nil
}

func (r *sessionRepository) GetByUserID(ctx context.Context, userID int) ([]*domain.Session, error) {
	var sessions []*domain.Session
	query := `
		SELECT * FROM sessions
		WHERE user_id = $1 AND expires_at > CURRENT_TIMESTAMP
		ORDER BY created_at DESC
	`
	err := r.db.SelectContext(ctx, &sessions, query, userID)
	return sessions, err
}

func (r *sessionRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM sessions WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrSessionNotFound
	}
	return nil
}

func (r *sessionRepository) DeleteByToken(ctx context.Context, token string) error {
	query := `DELETE FROM sessions WHERE token = $1`
	result, err := r.db.ExecContext(ctx, query, token)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrSessionNotFound
	}
	return nil
}

func (r *sessionRepository) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM sessions WHERE expires_at < CURRENT_TIMESTAMP`
	_, err := r.db.ExecContext(ctx, query)
	return err
}

func (r *sessionRepository) DeleteByUserID(ctx context.Context, userID int) error {
	query := `DELETE FROM sessions WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}
