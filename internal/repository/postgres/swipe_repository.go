package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gdugdh24/mpit2026-backend/internal/domain"
	"github.com/gdugdh24/mpit2026-backend/internal/repository"
	"github.com/jmoiron/sqlx"
)

type swipeRepository struct {
	db *sqlx.DB
}

func NewSwipeRepository(db *sqlx.DB) repository.SwipeRepository {
	return &swipeRepository{db: db}
}

func (r *swipeRepository) Create(ctx context.Context, swipe *domain.Swipe) error {
	query := `
		INSERT INTO swipes (swiper_id, swiped_id, is_like)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	return r.db.QueryRowContext(
		ctx, query,
		swipe.SwiperID, swipe.SwipedID, swipe.IsLike,
	).Scan(&swipe.ID, &swipe.CreatedAt)
}

func (r *swipeRepository) GetByID(ctx context.Context, id int) (*domain.Swipe, error) {
	var swipe domain.Swipe
	query := `SELECT * FROM swipes WHERE id = $1`
	err := r.db.GetContext(ctx, &swipe, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrSwipeAlreadyExists
		}
		return nil, err
	}
	return &swipe, nil
}

func (r *swipeRepository) GetByUsers(ctx context.Context, swiperID, swipedID int) (*domain.Swipe, error) {
	var swipe domain.Swipe
	query := `SELECT * FROM swipes WHERE swiper_id = $1 AND swiped_id = $2`
	err := r.db.GetContext(ctx, &swipe, query, swiperID, swipedID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &swipe, nil
}

func (r *swipeRepository) GetUserSwipes(ctx context.Context, userID int, limit, offset int) ([]*domain.Swipe, error) {
	var swipes []*domain.Swipe
	query := `
		SELECT * FROM swipes
		WHERE swiper_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	err := r.db.SelectContext(ctx, &swipes, query, userID, limit, offset)
	return swipes, err
}

func (r *swipeRepository) GetLikesReceived(ctx context.Context, userID int, limit, offset int) ([]*domain.Swipe, error) {
	var swipes []*domain.Swipe
	query := `
		SELECT * FROM swipes
		WHERE swiped_id = $1 AND is_like = true
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	err := r.db.SelectContext(ctx, &swipes, query, userID, limit, offset)
	return swipes, err
}

func (r *swipeRepository) CheckMutualLike(ctx context.Context, user1ID, user2ID int) (bool, error) {
	var count int
	query := `
		SELECT COUNT(*) FROM swipes
		WHERE ((swiper_id = $1 AND swiped_id = $2) OR (swiper_id = $2 AND swiped_id = $1))
		AND is_like = true
	`
	err := r.db.GetContext(ctx, &count, query, user1ID, user2ID)
	if err != nil {
		return false, err
	}
	return count == 2, nil
}
