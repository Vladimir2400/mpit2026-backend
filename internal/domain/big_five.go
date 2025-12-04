package domain

import "time"

type BigFiveResult struct {
	ID                int       `json:"id" db:"id"`
	UserID            int       `json:"user_id" db:"user_id"`
	Openness          float64   `json:"openness" db:"openness"`
	Conscientiousness float64   `json:"conscientiousness" db:"conscientiousness"`
	Extraversion      float64   `json:"extraversion" db:"extraversion"`
	Agreeableness     float64   `json:"agreeableness" db:"agreeableness"`
	Neuroticism       float64   `json:"neuroticism" db:"neuroticism"`
	CompletedAt       time.Time `json:"completed_at" db:"completed_at"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}
