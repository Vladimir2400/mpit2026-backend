package profile

import (
	"context"
	"fmt"

	"github.com/gdugdh24/mpit2026-backend/internal/domain"
	"github.com/gdugdh24/mpit2026-backend/internal/repository"
)

type ProfileUseCase struct {
	profileRepo repository.ProfileRepository
	userRepo    repository.UserRepository
}

func NewProfileUseCase(
	profileRepo repository.ProfileRepository,
	userRepo repository.UserRepository,
) *ProfileUseCase {
	return &ProfileUseCase{
		profileRepo: profileRepo,
		userRepo:    userRepo,
	}
}

// CreateProfileRequest represents profile creation request
type CreateProfileRequest struct {
	DisplayName       string   `json:"display_name" binding:"required,min=2,max=100"`
	Bio               *string  `json:"bio" binding:"omitempty,max=500"`
	City              *string  `json:"city" binding:"omitempty,max=100"`
	Interests         []string `json:"interests" binding:"omitempty,max=10"`
	PrefMinAge        *int     `json:"pref_min_age" binding:"omitempty,min=18,max=100"`
	PrefMaxAge        *int     `json:"pref_max_age" binding:"omitempty,min=18,max=100"`
	PrefMaxDistanceKm *int     `json:"pref_max_distance_km" binding:"omitempty,min=1,max=1000"`
}

// UpdateProfileRequest represents profile update request
type UpdateProfileRequest struct {
	DisplayName       *string   `json:"display_name" binding:"omitempty,min=2,max=100"`
	Bio               *string   `json:"bio" binding:"omitempty,max=500"`
	City              *string   `json:"city" binding:"omitempty,max=100"`
	Interests         *[]string `json:"interests" binding:"omitempty,max=10"`
	LocationLat       *float64  `json:"location_lat" binding:"omitempty,min=-90,max=90"`
	LocationLon       *float64  `json:"location_lon" binding:"omitempty,min=-180,max=180"`
	PrefMinAge        *int      `json:"pref_min_age" binding:"omitempty,min=18,max=100"`
	PrefMaxAge        *int      `json:"pref_max_age" binding:"omitempty,min=18,max=100"`
	PrefMaxDistanceKm *int      `json:"pref_max_distance_km" binding:"omitempty,min=1,max=1000"`
}

// ProfileResponse represents profile response with additional info
type ProfileResponse struct {
	*domain.Profile
	Age        int      `json:"age,omitempty"`
	DistanceKm *float64 `json:"distance_km,omitempty"`
}

// GetMyProfile returns current user's profile
func (uc *ProfileUseCase) GetMyProfile(ctx context.Context, userID int) (*domain.Profile, error) {
	profile, err := uc.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

// GetProfileByUserID returns profile by user ID with calculated age and distance
func (uc *ProfileUseCase) GetProfileByUserID(ctx context.Context, targetUserID int, currentUserID *int) (*ProfileResponse, error) {
	profile, err := uc.profileRepo.GetByUserID(ctx, targetUserID)
	if err != nil {
		return nil, err
	}

	user, err := uc.userRepo.GetByID(ctx, targetUserID)
	if err != nil {
		return nil, err
	}

	response := &ProfileResponse{
		Profile: profile,
		Age:     user.Age(),
	}

	// Calculate distance if current user location is available
	if currentUserID != nil {
		currentProfile, err := uc.profileRepo.GetByUserID(ctx, *currentUserID)
		if err == nil && currentProfile.LocationLat != nil && currentProfile.LocationLon != nil &&
			profile.LocationLat != nil && profile.LocationLon != nil {
			distance := calculateDistance(
				*currentProfile.LocationLat, *currentProfile.LocationLon,
				*profile.LocationLat, *profile.LocationLon,
			)
			response.DistanceKm = &distance
		}
	}

	return response, nil
}

// CreateProfile creates a new profile (onboarding)
func (uc *ProfileUseCase) CreateProfile(ctx context.Context, userID int, req *CreateProfileRequest) (*domain.Profile, error) {
	// Check if profile already exists
	existingProfile, err := uc.profileRepo.GetByUserID(ctx, userID)
	if err == nil && existingProfile != nil {
		return nil, domain.ErrProfileAlreadyExists
	}

	profile := &domain.Profile{
		UserID:               userID,
		DisplayName:          req.DisplayName,
		Bio:                  req.Bio,
		City:                 req.City,
		Interests:            req.Interests,
		PrefMinAge:           req.PrefMinAge,
		PrefMaxAge:           req.PrefMaxAge,
		PrefMaxDistanceKm:    req.PrefMaxDistanceKm,
		IsOnboardingComplete: true,
	}

	if err := uc.profileRepo.Create(ctx, profile); err != nil {
		return nil, fmt.Errorf("failed to create profile: %w", err)
	}

	return profile, nil
}

// UpdateProfile updates user profile
func (uc *ProfileUseCase) UpdateProfile(ctx context.Context, userID int, req *UpdateProfileRequest) (*domain.Profile, error) {
	profile, err := uc.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.DisplayName != nil {
		profile.DisplayName = *req.DisplayName
	}
	if req.Bio != nil {
		profile.Bio = req.Bio
	}
	if req.City != nil {
		profile.City = req.City
	}
	if req.Interests != nil {
		profile.Interests = *req.Interests
	}
	if req.LocationLat != nil {
		profile.LocationLat = req.LocationLat
	}
	if req.LocationLon != nil {
		profile.LocationLon = req.LocationLon
	}
	if req.PrefMinAge != nil {
		profile.PrefMinAge = req.PrefMinAge
	}
	if req.PrefMaxAge != nil {
		profile.PrefMaxAge = req.PrefMaxAge
	}
	if req.PrefMaxDistanceKm != nil {
		profile.PrefMaxDistanceKm = req.PrefMaxDistanceKm
	}

	if err := uc.profileRepo.Update(ctx, profile); err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	return profile, nil
}

// calculateDistance calculates distance between two points using Haversine formula
func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371 // km

	// Convert to radians
	lat1Rad := lat1 * 3.14159265359 / 180
	lat2Rad := lat2 * 3.14159265359 / 180
	deltaLat := (lat2 - lat1) * 3.14159265359 / 180
	deltaLon := (lon2 - lon1) * 3.14159265359 / 180

	// Haversine formula
	a := 0.5 - 0.5*cosine(deltaLat) + cosine(lat1Rad)*cosine(lat2Rad)*(1-cosine(deltaLon))/2

	return earthRadius * 2 * asin(sqrt(a))
}

// Helper math functions
func sqrt(x float64) float64 {
	if x < 0 {
		return 0
	}
	z := 1.0
	for i := 0; i < 10; i++ {
		z = z - (z*z-x)/(2*z)
	}
	return z
}

func cosine(x float64) float64 {
	// Taylor series approximation
	result := 1.0
	term := 1.0
	for i := 1; i <= 10; i++ {
		term *= -x * x / float64((2*i-1)*(2*i))
		result += term
	}
	return result
}

func asin(x float64) float64 {
	// Taylor series approximation
	if x < -1 || x > 1 {
		return 0
	}
	result := x
	term := x
	for i := 1; i <= 10; i++ {
		term *= x * x * float64((2*i-1)*(2*i-1)) / float64((2*i)*(2*i+1))
		result += term
	}
	return result
}
