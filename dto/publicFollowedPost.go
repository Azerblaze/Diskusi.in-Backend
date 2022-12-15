package dto

import (
	"gorm.io/gorm"
)

type PublicFollowedPost struct {
	gorm.Model
	User FollowedPostUser `json:"user" form:"user"`
	Post FollowedPost     `json:"post" form:"post"`
}
type FollowedPost struct {
	PostID    int    `json:"postId" form:"postId"`
	PostTopic string `json:"postTopic" form:"postTopic"`
	Title     string `json:"title" form:"title"`
	Body      string `json:"body" form:"body"`
}
type FollowedPostUser struct {
	UserID   int    `json:"userId" form:"userId"`
	Photo    string `json:"photo" form:"photo"`
	Username string `json:"username" form:"username"`
}
