package domain

import "time"

type UserEmbedding struct {
	ID                int        `json:"id" db:"id"`
	UserID            int        `json:"user_id" db:"user_id"`
	MusicVector       []float32  `json:"music_vector" db:"music_vector"`
	GroupsVector      []float32  `json:"groups_vector" db:"groups_vector"`
	PostsVector       []float32  `json:"posts_vector" db:"posts_vector"`
	CombinedVector    []float32  `json:"combined_vector" db:"combined_vector"`
	VKDataFetchedAt   *time.Time `json:"vk_data_fetched_at" db:"vk_data_fetched_at"`
	VectorsUpdatedAt  *time.Time `json:"vectors_updated_at" db:"vectors_updated_at"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
}
