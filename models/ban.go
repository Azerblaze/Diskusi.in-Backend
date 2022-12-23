package models

import (
	"gorm.io/gorm"
)

type Ban struct {
	gorm.Model
	UserID    int `json:"userId" form:"userId"`
	TopicID   int `json:"topicId" form:"topicId"`
	Ban_Until int `json:"banUntil" form:"banUntil"`

	User  User  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Topic Topic `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (Ban) TableName() string {
	return "bans"
}
