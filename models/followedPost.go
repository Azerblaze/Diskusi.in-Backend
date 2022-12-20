package models

import (
	"gorm.io/gorm"
)

type FollowedPost struct {
	gorm.Model
	UserID int `json:"userId" form:"userId"`
	PostID int `json:"postId" form:"postId"`

	User User `json:"user" form:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Post Post `json:"post" form:"post" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (FollowedPost) TableName() string {
	return "followedPosts"
}
