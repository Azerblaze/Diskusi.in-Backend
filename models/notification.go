package models

import (
	"gorm.io/gorm"
)

type Notification struct {
	gorm.Model
	UserID      int    `json:"userId" form:"userId"`
	PostID      int    `json:"postId" form:"postId"`
	IsSee       bool   `json:"isSee" form:"isSee"`
	Category    string `json:"category" form:"category"`
	Description string `json:"description" form:"description"`

	User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Post Post `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (Notification) TableName() string {
	return "notifications"
}
