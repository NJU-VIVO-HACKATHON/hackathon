package model

import "gorm.io/gorm"

type PostTag struct {
	gorm.Model
	Pid int64
	Tid int64
}
