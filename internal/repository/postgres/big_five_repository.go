package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gdugdh24/mpit2026-backend/internal/domain"
	"github.com/gdugdh24/mpit2026-backend/internal/repository"
	"github.com/jmoiron/sqlx"
)

type bigFiveRepository struct {
	db *sqlx.DB
}

func NewBigFiveRepository(db *sqlx.DB) repository.BigFiveRepository {
	return &bigFiveRepository{db: db}
}

func (r *bigFiveRepository) Create(ctx context.Context, result *domain.BigFiveResult) error {
	query := `
		INSERT INTO big_five_results (
			user_id, openness, conscientiousness, extraversion,
			agreeableness, neuroticism, completed_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(
		ctx, query,
		result.UserID, result.Openness, result.Conscientiousness,
		result.Extraversion, result.Agreeableness, result.Neuroticism,
		result.CompletedAt,
	).Scan(&result.ID, &result.CreatedAt, &result.UpdatedAt)
}

func (r *bigFiveRepository) GetByID(ctx context.Context, id int) (*domain.BigFiveResult, error) {
	var result domain.BigFiveResult
	query := `SELECT * FROM big_five_results WHERE id = $1`
	err := r.db.GetContext(ctx, &result, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("big five result not found")
		}
		return nil, err
	}
	return &result, nil
}

func (r *bigFiveRepository) GetByUserID(ctx context.Context, userID int) (*domain.BigFiveResult, error) {
	var result domain.BigFiveResult
	query := `SELECT * FROM big_five_results WHERE user_id = $1`
	err := r.db.GetContext(ctx, &result, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("big five result not found")
		}
		return nil, err
	}
	return &result, nil
}

func (r *bigFiveRepository) Update(ctx context.Context, result *domain.BigFiveResult) error {
	query := `
		UPDATE big_five_results
		SET openness = $1, conscientiousness = $2, extraversion = $3,
		    agreeableness = $4, neuroticism = $5, completed_at = $6,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $7
		RETURNING updated_at
	`
	return r.db.QueryRowContext(
		ctx, query,
		result.Openness, result.Conscientiousness, result.Extraversion,
		result.Agreeableness, result.Neuroticism, result.CompletedAt,
		result.ID,
	).Scan(&result.UpdatedAt)
}

func (r *bigFiveRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM big_five_results WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("big five result not found")
	}
	return nil
}
