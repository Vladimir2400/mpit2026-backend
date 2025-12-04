package domain

import "time"

type Swipe struct {
	ID        int       `json:"id" db:"id"`
	SwiperID  int       `json:"swiper_id" db:"swiper_id"`
	SwipedID  int       `json:"swiped_id" db:"swiped_id"`
	IsLike    bool      `json:"is_like" db:"is_like"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
