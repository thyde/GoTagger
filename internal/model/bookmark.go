package model

import "time"

type Bookmark struct {
	ID        int64     `json:"id" db:"id"`
	URL       string    `json:"url" db:"url"`
	Title     string    `json:"title" db:"title"`
	Tags      []string  `json:"tags" db:"tags"`
	Favorite  bool      `json:"favorite" db:"favorite"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
