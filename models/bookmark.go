package models

import (
	"gorm.io/gorm"
)

type Bookmark struct {
	gorm.Model
	UserID int `json:"userId" form:"userId"`
	PostID int `json:"postId" form:"postId"`

	// User User
	Post Post `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (Bookmark) TableName() string {
	return "bookmarks"
}
