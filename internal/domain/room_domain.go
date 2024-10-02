package domain

import "time"

// Room
type Room struct {
	ID        string `gorm:"unique"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Rooms []Room
