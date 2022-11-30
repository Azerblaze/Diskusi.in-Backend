package repositories

import (
	"discusiin/models"
)

type IDatabase interface {
	SaveNewUser(models.User) error
	Login(username string, password string) (models.User, error)
	GetUserByUsername(username string) (models.User, error)

	GetAllTopics() ([]models.Topic, error)
	GetTopicByName(topicName string) (models.Topic, error)
	GetTopicByID(topicID int) (models.Topic, error)
	SaveNewTopic(models.Topic) error
	SaveTopic(models.Topic) error
	RemoveTopic(topicID int) error

	SaveNewPost(post models.Post) error
	GetAllPostByTopic(topidID int) ([]models.Post, error)
	GetPostById(postID int) (models.Post, error)
	SavePost(post models.Post) error
	DeletePost(postID int) error
	GetPostByIdWithAll(postID int) (models.Post, error)

	SaveNewComment(comment models.Comment) error
	GetAllCommentByPost(postID int) ([]models.Comment, error)
	GetCommentById(commendID int) (models.Comment, error)
	SaveComment(comment models.Comment) error
	DeleteComment(commentID int) error

	SaveNewReply(reply models.Reply) error
	GetAllReplyByComment(commentId int) ([]models.Reply, error)
	GetReplyById(re int) (models.Reply, error)
	SaveReply(reply models.Reply) error
	DeleteReply(re int) error
}
