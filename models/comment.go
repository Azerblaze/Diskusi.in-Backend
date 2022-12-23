package models

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	UserID     int    `json:"userId" form:"userId"`
	PostID     int    `json:"postId" form:"postId"`
	Body       string `json:"body" form:"body"`
	IsFollowed bool   `json:"isFollowed" form:"isFollowed"`

	User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Post Post `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (Comment) TableName() string {
	return "comments"
}
