package domain

import (
	"errors"
	"regexp"
	"time"
)

const (
	nameLengthMax = 100
	nameLengthMin = 1
	passLengthMax = 100
	passLengthMin = 8
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

	if len(u.Password) < passLengthMin || len(u.Password) > passLengthMax {
		return errors.New("password の値は8文字以上100文字以内にしてください")
	}

	passPattern := `^[a-zA-Z\d!@#$%^&*()_+\-=\[\]]{8,100}$`
	matched, err := regexp.MatchString(passPattern, u.Password)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("パスワードは半角英数字をそれぞれ1種類以上含み、8文字以上100文字以内である必要があります。")
	}

	hasLetter := false
	hasDigit := false

	for _, char := range u.Password {
		if char >= 'A' && char <= 'Z' || char >= 'a' && char <= 'z' {
			hasLetter = true
		}
		if char >= '0' && char <= '9' {
			hasDigit = true
		}
	}

	if !(hasLetter && hasDigit) {
		return errors.New("パスワードは半角英数字をそれぞれ1種類以上含み、8文字以上100文字以内である必要があります。")
	}

	return nil
}
