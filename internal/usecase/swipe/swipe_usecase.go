package swipe

import (
	"context"
	"fmt"

	"github.com/gdugdh24/mpit2026-backend/internal/domain"
	"github.com/gdugdh24/mpit2026-backend/internal/repository"
)

type SwipeUseCase struct {
	swipeRepo   repository.SwipeRepository
	matchRepo   repository.MatchRepository
	profileRepo repository.ProfileRepository
	userRepo    repository.UserRepository
}

func NewSwipeUseCase(
	swipeRepo repository.SwipeRepository,
	matchRepo repository.MatchRepository,
	profileRepo repository.ProfileRepository,
	userRepo repository.UserRepository,
) *SwipeUseCase {
	return &SwipeUseCase{
		swipeRepo:   swipeRepo,
		matchRepo:   matchRepo,
		profileRepo: profileRepo,
		userRepo:    userRepo,
	}
}

// SwipeRequest represents a swipe action
type SwipeRequest struct {
	SwipedUserID int  `json:"swiped_user_id" binding:"required"`
	IsLike       bool `json:"is_like"`
}

// SwipeResponse represents swipe result
type SwipeResponse struct {
	IsMatch     bool                `json:"is_match"`
	Swipe       *domain.Swipe       `json:"swipe,omitempty"`
	Match       *domain.Match       `json:"match,omitempty"`
	MatchedUser *MatchedUserProfile `json:"matched_user,omitempty"`
}

// MatchedUserProfile represents matched user info
type MatchedUserProfile struct {
	ID          int      `json:"id"`
	DisplayName string   `json:"display_name"`
	Bio         *string  `json:"bio"`
	City        *string  `json:"city"`
	Age         int      `json:"age"`
	DistanceKm  *float64 `json:"distance_km"`
}

// LikeReceivedResponse represents a like received
type LikeReceivedResponse struct {
	SwipeID   int                 `json:"swipe_id"`
	User      *MatchedUserProfile `json:"user"`
	CreatedAt string              `json:"created_at"`
}

// CreateSwipe creates a new swipe and checks for match
func (uc *SwipeUseCase) CreateSwipe(ctx context.Context, swiperID int, req *SwipeRequest) (*SwipeResponse, error) {
	// Validate: can't swipe yourself
	if swiperID == req.SwipedUserID {
		return nil, domain.ErrCannotSwipeSelf
	}

	// Check if already swiped
	existingSwipe, err := uc.swipeRepo.GetByUsers(ctx, swiperID, req.SwipedUserID)
	if err == nil && existingSwipe != nil {
		return nil, domain.ErrSwipeAlreadyExists
	}

	// Create swipe
	swipe := &domain.Swipe{
		SwiperID: swiperID,
		SwipedID: req.SwipedUserID,
		IsLike:   req.IsLike,
	}

	if err := uc.swipeRepo.Create(ctx, swipe); err != nil {
		return nil, fmt.Errorf("failed to create swipe: %w", err)
	}

	response := &SwipeResponse{
		IsMatch: false,
		Swipe:   swipe,
	}

	// If it's a like, check for mutual like (match)
	if req.IsLike {
		isMutual, err := uc.swipeRepo.CheckMutualLike(ctx, swiperID, req.SwipedUserID)
		if err != nil {
			return response, nil // Return swipe even if match check fails
		}

		if isMutual {
			// Create match
			match, err := uc.createMatch(ctx, swiperID, req.SwipedUserID)
			if err != nil {
				return response, nil // Return swipe even if match creation fails
			}

			// Get matched user profile
			matchedUser, err := uc.getMatchedUserProfile(ctx, req.SwipedUserID)
			if err == nil {
				response.IsMatch = true
				response.Match = match
				response.MatchedUser = matchedUser
			}
		}
	}

	return response, nil
}

// createMatch creates a match between two users
func (uc *SwipeUseCase) createMatch(ctx context.Context, user1ID, user2ID int) (*domain.Match, error) {
	// Ensure user1_id < user2_id for database constraint
	if user1ID > user2ID {
		user1ID, user2ID = user2ID, user1ID
	}

	// Check if match already exists
	existingMatch, err := uc.matchRepo.GetByUsers(ctx, user1ID, user2ID)
	if err == nil && existingMatch != nil {
		return existingMatch, nil
	}

	match := &domain.Match{
		User1ID:  user1ID,
		User2ID:  user2ID,
		IsActive: true,
	}

	if err := uc.matchRepo.Create(ctx, match); err != nil {
		return nil, err
	}

	return match, nil
}

// getMatchedUserProfile gets basic profile info for matched user
func (uc *SwipeUseCase) getMatchedUserProfile(ctx context.Context, userID int) (*MatchedUserProfile, error) {
	profile, err := uc.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &MatchedUserProfile{
		ID:          profile.ID,
		DisplayName: profile.DisplayName,
		Bio:         profile.Bio,
		City:        profile.City,
		Age:         user.Age(),
	}, nil
}

// GetLikesReceived returns list of users who liked current user
func (uc *SwipeUseCase) GetLikesReceived(ctx context.Context, userID int, limit, offset int) ([]*LikeReceivedResponse, int, error) {
	likes, err := uc.swipeRepo.GetLikesReceived(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get likes received: %w", err)
	}

	responses := make([]*LikeReceivedResponse, 0, len(likes))
	for _, like := range likes {
		// Get user profile
		profile, err := uc.profileRepo.GetByUserID(ctx, like.SwiperID)
		if err != nil {
			continue
		}

		user, err := uc.userRepo.GetByID(ctx, like.SwiperID)
		if err != nil {
			continue
		}

		// Calculate distance if location available
		var distanceKm *float64
		currentProfile, err := uc.profileRepo.GetByUserID(ctx, userID)
		if err == nil && currentProfile.LocationLat != nil && currentProfile.LocationLon != nil &&
			profile.LocationLat != nil && profile.LocationLon != nil {
			distance := calculateDistance(
				*currentProfile.LocationLat, *currentProfile.LocationLon,
				*profile.LocationLat, *profile.LocationLon,
			)
			distanceKm = &distance
		}

		responses = append(responses, &LikeReceivedResponse{
			SwipeID: like.ID,
			User: &MatchedUserProfile{
				ID:          profile.ID,
				DisplayName: profile.DisplayName,
				Bio:         profile.Bio,
				City:        profile.City,
				Age:         user.Age(),
				DistanceKm:  distanceKm,
			},
			CreatedAt: like.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return responses, len(likes), nil
}

// calculateDistance calculates distance between two points
func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	// Simple Haversine formula implementation
	// This is simplified - use math package for production
	const earthRadius = 6371.0                   // km
	dLat := (lat2 - lat1) * 0.017453292519943295 // to radians
	dLon := (lon2 - lon1) * 0.017453292519943295
	a := 0.5 - 0.5*cosApprox(dLat) + cosApprox(lat1*0.017453292519943295)*cosApprox(lat2*0.017453292519943295)*(1-cosApprox(dLon))/2
	return earthRadius * 2 * asinApprox(sqrtApprox(a))
}

func cosApprox(x float64) float64 {
	// Taylor series approximation
	x2 := x * x
	return 1 - x2/2 + x2*x2/24
}

func sqrtApprox(x float64) float64 {
	if x <= 0 {
		return 0
	}
	z := x
	for i := 0; i < 10; i++ {
		z = (z + x/z) / 2
	}
	return z
}

func asinApprox(x float64) float64 {
	if x < -1 || x > 1 {
		return 0
	}
	return x + x*x*x/6 + 3*x*x*x*x*x/40
}
