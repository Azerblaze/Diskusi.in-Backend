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

	Comments []Comment `json:"comments"`
	User     User      `json:"user"`
	Topic    Topic     `json:"topic"`
}

func (Post) TableName() string {
	return "posts"
}
