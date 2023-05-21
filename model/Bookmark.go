package model

import "gorm.io/gorm"

type Enum int

type Bookmark struct {
	gorm.Model
	Uid      *int64
	Pid      *int64
	Like     bool
	Favorite bool
}
