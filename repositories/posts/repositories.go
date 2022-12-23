package posts

import (
	"discusiin/models"
)

type IPostDatabase interface {
	SaveNewPost(post models.Post) error
	GetAllPostByTopic(topidID int, page int, search string) ([]models.Post, error)
	GetAllPostByTopicByLike(topicID int, page int) ([]models.Post, error)
	GetPostById(postID int) (models.Post, error)
	GetRecentPost(page int, search string) ([]models.Post, error)
	GetPostByUserId(userId int, page int) ([]models.Post, error)
	GetAllPostByLike(page int) ([]models.Post, error)
	SavePost(post models.Post) error
	DeletePostByPostID(postID int) error
	DeletePostByUserID(userID int) error
	GetPostByIdWithAll(postID int) (models.Post, error)
	CountPostLike(postID int) (int, error)
	CountPostComment(postID int) (int, error)
	CountPostDislike(postID int) (int, error)
	CountPostByTopicID(topicId int) (int, error)
	CountPostByUserID(userId int) (int, error)
	CountAllPost() (int, error)
	CountNumberOfPostByTopicName(topicName string) (int, error)
}
