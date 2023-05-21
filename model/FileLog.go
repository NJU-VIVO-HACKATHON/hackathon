package model

import "gorm.io/gorm"

type UploadFileLog struct {
	gorm.Model
	FileName  string
	UserAgent string
	FileType  string
	FileSize  int64
	SavePath  string
}
