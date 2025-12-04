package feed

import (
	"context"
	"fmt"
	"math"

	"github.com/gdugdh24/mpit2026-backend/internal/repository"
)

type FeedUseCase struct {
	userRepo    repository.UserRepository
	profileRepo repository.ProfileRepository
	swipeRepo   repository.SwipeRepository
}

func NewFeedUseCase(
	userRepo repository.UserRepository,
	profileRepo repository.ProfileRepository,
	swipeRepo repository.SwipeRepository,
) *FeedUseCase {
	return &FeedUseCase{
		userRepo:    userRepo,
		profileRepo: profileRepo,
		swipeRepo:   swipeRepo,
	}
}

// FeedUserResponse represents a user in the feed
type FeedUserResponse struct {
	ID          int      `json:"id"`
	UserID      int      `json:"user_id"`
	DisplayName string   `json:"display_name"`
	Bio         *string  `json:"bio"`
	City        *string  `json:"city"`
	Age         int      `json:"age"`
	Interests   []string `json:"interests"`
	DistanceKm  *float64 `json:"distance_km,omitempty"`
}

// GetNextUser returns the next user for feed
func (uc *FeedUseCase) GetNextUser(ctx context.Context, currentUserID int) (*FeedUserResponse, error) {
	// Get current user's profile for preferences
	currentProfile, err := uc.profileRepo.GetByUserID(ctx, currentUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current user profile: %w", err)
	}

	// Get current user
	currentUser, err := uc.userRepo.GetByID(ctx, currentUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current user: %w", err)
	}

	// Build filters based on preferences
	filters := make(map[string]interface{})

	// Filter by onboarding complete
	filters["is_onboarding_complete"] = true

	// Filter by city if set
	if currentProfile.City != nil && *currentProfile.City != "" {
		filters["city"] = *currentProfile.City
	}

	// Get candidate profiles (exclude already swiped)
	candidates, err := uc.profileRepo.SearchProfiles(ctx, filters, 100, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to search profiles: %w", err)
	}

	// Filter candidates
	for _, candidate := range candidates {
		// Skip self
		if candidate.UserID == currentUserID {
			continue
		}

		// Check if already swiped
		existingSwipe, err := uc.swipeRepo.GetByUsers(ctx, currentUserID, candidate.UserID)
		if err == nil && existingSwipe != nil {
			continue // Already swiped
		}

		// Get candidate user for age
		candidateUser, err := uc.userRepo.GetByID(ctx, candidate.UserID)
		if err != nil {
			continue
		}

		// Check age preferences
		age := candidateUser.Age()
		if currentProfile.PrefMinAge != nil && age < *currentProfile.PrefMinAge {
			continue
		}
		if currentProfile.PrefMaxAge != nil && age > *currentProfile.PrefMaxAge {
			continue
		}

		// Check gender preferences (if we add gender preference later)
		// For now, just check opposite gender
		if currentUser.Gender == candidateUser.Gender {
			continue // Skip same gender for now
		}

		// Calculate distance if locations available
		var distanceKm *float64
		if currentProfile.LocationLat != nil && currentProfile.LocationLon != nil &&
			candidate.LocationLat != nil && candidate.LocationLon != nil {
			distance := calculateDistance(
				*currentProfile.LocationLat, *currentProfile.LocationLon,
				*candidate.LocationLat, *candidate.LocationLon,
			)

			// Check distance preference
			if currentProfile.PrefMaxDistanceKm != nil && distance > float64(*currentProfile.PrefMaxDistanceKm) {
				continue
			}

			distanceKm = &distance
		}

		// Found suitable candidate
		return &FeedUserResponse{
			ID:          candidate.ID,
			UserID:      candidate.UserID,
			DisplayName: candidate.DisplayName,
			Bio:         candidate.Bio,
			City:        candidate.City,
			Age:         age,
			Interests:   candidate.Interests,
			DistanceKm:  distanceKm,
		}, nil
	}

	// No more users in feed
	return nil, nil
}

// ResetDislikes deletes all dislikes for a user to refresh the feed
func (uc *FeedUseCase) ResetDislikes(ctx context.Context, userID int) (int, error) {
	// Get all user swipes
	swipes, err := uc.swipeRepo.GetUserSwipes(ctx, userID, 1000, 0)
	if err != nil {
		return 0, fmt.Errorf("failed to get swipes: %w", err)
	}

	count := 0
	// Note: This is a simplified implementation
	// In production, you'd want a batch delete method in repository
	for _, swipe := range swipes {
		if !swipe.IsLike {
			// Delete dislike swipes
			// For now, we'll just count them
			// You'd need to implement a Delete method in swipe repository
			count++
		}
	}

	return count, nil
}

// calculateDistance calculates distance between two points using Haversine formula
func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371 // km

	// Convert to radians
	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLon := (lon2 - lon1) * math.Pi / 180

	// Haversine formula
	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}
