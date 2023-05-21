package model

import "gorm.io/gorm"

type LoveTag struct {
	gorm.Model
	Uid   int64
	Tid   int64
	Level int64
}
