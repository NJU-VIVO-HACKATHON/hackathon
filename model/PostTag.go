package model

import "gorm.io/gorm"

type PostTag struct {
	gorm.Model
	Pid string
	Tid string
}
