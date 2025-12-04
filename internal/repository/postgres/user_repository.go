package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/gdugdh24/mpit2026-backend/internal/domain"
	"github.com/gdugdh24/mpit2026-backend/internal/repository"
	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (vk_id, vk_access_token, vk_token_expires_at, gender, birth_date, is_verified, is_online)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(
		ctx, query,
		user.VKID, user.VKAccessToken, user.VKTokenExpiresAt,
		user.Gender, user.BirthDate, user.IsVerified, user.IsOnline,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	var user domain.User
	query := `SELECT * FROM users WHERE id = $1`
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByVKID(ctx context.Context, vkID int) (*domain.User, error) {
	var user domain.User
	query := `SELECT * FROM users WHERE vk_id = $1`
	err := r.db.GetContext(ctx, &user, query, vkID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET vk_access_token = $1, vk_token_expires_at = $2, gender = $3,
		    birth_date = $4, is_verified = $5, is_online = $6,
		    last_online_at = $7, updated_at = CURRENT_TIMESTAMP
		WHERE id = $8
		RETURNING updated_at
	`
	return r.db.QueryRowContext(
		ctx, query,
		user.VKAccessToken, user.VKTokenExpiresAt, user.Gender,
		user.BirthDate, user.IsVerified, user.IsOnline,
		user.LastOnlineAt, user.ID,
	).Scan(&user.UpdatedAt)
}

func (r *userRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (r *userRepository) UpdateOnlineStatus(ctx context.Context, userID int, isOnline bool) error {
	now := time.Now()
	query := `
		UPDATE users
		SET is_online = $1, last_online_at = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`
	result, err := r.db.ExecContext(ctx, query, isOnline, now, userID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (r *userRepository) GetOnlineUsers(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	var users []*domain.User
	query := `
		SELECT * FROM users
		WHERE is_online = true
		ORDER BY last_online_at DESC
		LIMIT $1 OFFSET $2
	`
	err := r.db.SelectContext(ctx, &users, query, limit, offset)
	return users, err
}
