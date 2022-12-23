package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title     string `json:"title" form:"title"`
	Photo     string `json:"photo" form:"photo"`
	Body      string `json:"body" form:"body"`
	UserID    int    `json:"userId" form:"userId"`
	TopicID   int    `json:"topicId" form:"topicId"`
	CreatedAt int    `json:"createdAt" form:"createdAt"`
	IsActive  bool   `json:"isActive" form:"isActive"`
	LikeCount int    `json:"likecount" form:"likecount"`

	Comments []Comment `json:"comments" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	User     User      `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Topic    Topic     `json:"topic" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (Post) TableName() string {
	return "posts"
}
