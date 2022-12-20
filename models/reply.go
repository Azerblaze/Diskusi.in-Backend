package models

import (
	"gorm.io/gorm"
)

type Reply struct {
	gorm.Model
	UserID    int    `json:"userId" form:"userId"`
	CommentID int    `json:"commentId" form:"commentId"`
	Body      string `json:"body" form:"body"`

	User    User    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Comment Comment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (Reply) TableName() string {
	return "replies"
}
