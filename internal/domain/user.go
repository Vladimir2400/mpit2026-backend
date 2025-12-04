package domain

import "time"

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

type User struct {
	ID                int        `json:"id" db:"id"`
	VKID              int        `json:"vk_id" db:"vk_id"`
	VKAccessToken     *string    `json:"-" db:"vk_access_token"`
	VKTokenExpiresAt  *time.Time `json:"-" db:"vk_token_expires_at"`
	Gender            Gender     `json:"gender" db:"gender"`
	BirthDate         time.Time  `json:"birth_date" db:"birth_date"`
	IsVerified        bool       `json:"is_verified" db:"is_verified"`
	IsOnline          bool       `json:"is_online" db:"is_online"`
	LastOnlineAt      *time.Time `json:"last_online_at" db:"last_online_at"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
}

func (u *User) Age() int {
	return int(time.Since(u.BirthDate).Hours() / 24 / 365.25)
}

func (u *User) IsAdult() bool {
	return u.Age() >= 18
}
