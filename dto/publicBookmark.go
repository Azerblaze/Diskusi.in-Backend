package dto

import (
	"gorm.io/gorm"
)

type PublicBookmark struct {
	gorm.Model
	User BookmarkUser `json:"user" form:"user"`
	Post BookmarkPost `json:"post" form:"post"`
}
type BookmarkPost struct {
	PostID    int    `json:"postId" form:"postId"`
	PostTopic string `json:"postTopic" form:"postTopic"`
	Title     string `json:"title" form:"title"`
	Body      string `json:"body" form:"body"`
}
type BookmarkUser struct {
	UserID   int    `json:"userId" form:"userId"`
	Photo    string `json:"photo" form:"photo"`
	Username string `json:"username" form:"username"`
}
