package dto

import "gorm.io/gorm"

type AdminComment struct {
	gorm.Model
	Body string      `json:"body" form:"body"`
	Post CommentPost `json:"post" form:"post"`
}

type CommentPost struct {
	PostID int    `json:"postId" form:"postId"`
	Title  string `json:"title" form:"title"`
	Body   string `json:"body" form:"body"`
}
