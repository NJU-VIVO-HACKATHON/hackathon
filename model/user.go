package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email        *string
	Sms          *string
	Nickname     *string
	Avatar       *string
	Introduction *string
}
