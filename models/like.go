package models

import (
	"gorm.io/gorm"
)

type Like struct {
	gorm.Model
	UserID    int  `json:"userId" form:"userId"`
	PostID    int  `json:"postId" form:"postId"`
	IsLike    bool `json:"isLike" form:"isLike"`
	IsDislike bool `json:"isDislike" form:"isDislike"`

	User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Post Post `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (Like) TableName() string {
	return "likes"
}
