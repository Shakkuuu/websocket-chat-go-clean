package domain

import "time"

// ユーザーの参加中Room情報
type ParticipatingRoom struct {
	ID        int `gorm:"unique"`
	RoomID    string
	IsMaster  bool
	UserID    string `gorm:"foreignKey:UserID;references:ID"`
	User      User
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ParticipatingRooms []ParticipatingRoom
