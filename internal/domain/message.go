package domain

import "time"

type Message struct {
	ID        int       `json:"id" db:"id"`
	MatchID   int       `json:"match_id" db:"match_id"`
	SenderID  int       `json:"sender_id" db:"sender_id"`
	Content   string    `json:"content" db:"content"`
	IsRead    bool      `json:"is_read" db:"is_read"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
