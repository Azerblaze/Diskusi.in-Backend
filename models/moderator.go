package models

import (
	"gorm.io/gorm"
)

type Moderator struct {
	gorm.Model
	UserID  int `json:"userId" form:"userId"`
	TopicID int `json:"topicId" form:"topicId"`

	User  User  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Topic Topic `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (Moderator) TableName() string {
	return "moderators"
}
