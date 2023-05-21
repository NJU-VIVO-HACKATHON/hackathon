package model

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Uid           *int64
	Title         *string
	Content       *string
	LikeCount     *int64
	FavoriteCount *int64
	Cover         *string
}
