package models

import "time"

type Item struct {
	ID         int     `json:"id" db:"id"`
	Name       string  `json:"name" db:"name"`
	Rarity     string  `json:"rarity" db:"rarity"`
	BaseChance float64 `json:"base_chance" db:"base_chance"`
	MinPity    int     `json:"min_pity" db:"min_pity"`
}

type DropEvent struct {
	ID           int       `json:"id" db:"id"`
	UserID       int       `json:"user_id" db:"user_id"`
	ItemID       int       `json:"item_id" db:"item_id"`
	DroppedAt    time.Time `json:"dropped_at" db:"dropped_at"`
	IsGuaranteed bool      `json:"is_guaranteed" db:"is_guaranteed"`
}
