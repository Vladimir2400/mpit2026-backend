package domain

import "time"

type Session struct {
	ID         int       `json:"id" db:"id"`
	UserID     int       `json:"user_id" db:"user_id"`
	Token      string    `json:"token" db:"token"`
	DeviceInfo *string   `json:"device_info" db:"device_info"`
	IPAddress  *string   `json:"ip_address" db:"ip_address"`
	ExpiresAt  time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
