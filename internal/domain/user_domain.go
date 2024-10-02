package domain

import (
	"errors"
	"time"
)

const (
	nameLengthMax = 100
	nameLengthMin = 1
	passLengthMin = 6
)

type Users []User

type User struct {
	ID        string `gorm:"unique"`
	Name      string `gorm:"unique"`
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) Validate() error {
	if u.Name == "" {
		return errors.New("name の値が不正です。")
	}

	if u.Password == "" {
		return errors.New("password の値が不正です。")
	}

	if len(u.Name) < nameLengthMin || len(u.Name) > nameLengthMax {
		return errors.New("name の値は1文字以上100文字以内にしてください")
	}

	if len(u.Password) < passLengthMin {
		return errors.New("password の値は6文字以上にしてください")
	}

	return nil
}
