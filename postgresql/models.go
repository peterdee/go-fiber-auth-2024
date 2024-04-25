package postgresql

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email     string `gorm:"UNIQUE"`
	FirstName string
	LastName  string

	Password          Password           `json:"-"`
	UsedRefreshTokens []UsedRefreshToken `json:"-"`
	UserSecret        UserSecret         `json:"-"`
}

type Password struct {
	gorm.Model
	Hash   string
	UserID uint
}

type UserSecret struct {
	gorm.Model
	Secret string
	UserID uint
}

type UsedRefreshToken struct {
	gorm.Model
	ExpiresAt int64
	Token     string
	UserID    uint
}
