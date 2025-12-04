package domain

import "errors"

var (
	// User errors
	ErrUserNotFound         = errors.New("user not found")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrInvalidVKSignature   = errors.New("invalid VK signature")
	ErrVKTokenExpired       = errors.New("VK token expired")

	// Profile errors
	ErrProfileNotFound      = errors.New("profile not found")
	ErrProfileAlreadyExists = errors.New("profile already exists")

	// Session errors
	ErrSessionNotFound      = errors.New("session not found")
	ErrSessionExpired       = errors.New("session expired")
	ErrInvalidToken         = errors.New("invalid token")

	// Swipe errors
	ErrSwipeAlreadyExists   = errors.New("swipe already exists")
	ErrCannotSwipeSelf      = errors.New("cannot swipe yourself")

	// Match errors
	ErrMatchNotFound        = errors.New("match not found")
	ErrNotMatched           = errors.New("users are not matched")
	ErrMatchAlreadyExists   = errors.New("match already exists")

	// Message errors
	ErrMessageNotFound      = errors.New("message not found")
	ErrUnauthorizedMessage  = errors.New("unauthorized to access message")

	// General errors
	ErrInvalidInput         = errors.New("invalid input")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrForbidden            = errors.New("forbidden")
	ErrInternalServer       = errors.New("internal server error")
)
