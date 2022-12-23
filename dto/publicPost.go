package dto

import (
	"gorm.io/gorm"
)

type PublicPost struct {
	gorm.Model
	Title     string    `json:"title" form:"title"`
	Photo     string    `json:"photo" form:"photo"`
	Body      string    `json:"body" form:"body"`
	CreatedAt int       `json:"createdAt" form:"createdAt"`
	IsActive  bool      `json:"isActive" form:"isActive"`
	User      PostUser  `json:"user" form:"user"`
	Topic     PostTopic `json:"topic" form:"topic"`
	Count     PostCount `json:"count" form:"count"`
}
type PostUser struct {
	UserID   int    `json:"userId" form:"userId"`
	Photo    string `json:"photo" form:"photo"`
	Username string `json:"username" form:"username"`
}
type PostTopic struct {
	TopicID   int    `json:"topicId" form:"topicId"`
	TopicName string `json:"topicName" form:"topicName"`
}
type PostCount struct {
	LikeCount    int `json:"like" form:"like"`
	CommentCount int `json:"comment" form:"comment"`
	DislikeCount int `json:"dislike" form:"dislike"`
}
