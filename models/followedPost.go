package models

import (
	"gorm.io/gorm"
)

type FollowedPost struct {
	gorm.Model
	UserID int `json:"userId" form:"userId"`
	PostID int `json:"postId" form:"postId"`

	User User `json:"user" form:"user"`
	Post Post `json:"post" form:"post"`
}

func (FollowedPost) TableName() string {
	return "followedPosts"
}
