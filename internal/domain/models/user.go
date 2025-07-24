package models

import "time"

type User struct {
	ID          int       `json:"id" db:"id"`
	Username    string    `json:"username" db:"username"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	PityCounter int       `json:"pity_counter" db:"pity_counter"`
}
