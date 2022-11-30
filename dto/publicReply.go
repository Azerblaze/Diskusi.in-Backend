package dto

import "gorm.io/gorm"

type PublicReply struct {
	gorm.Model
	UserID    int    `json:"userId" form:"userId"`
	CommentID int    `json:"commentId" form:"commentId"`
	Body      string `json:"body" form:"body"`
	Username  string `json:"username" form:"username"`
}
