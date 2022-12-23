package followedPosts

import (
	"discusiin/models"
)

type IFollowedPostDatabase interface {
	SaveFollowedPost(followedPost models.FollowedPost) error
	GetFollowedPost(userId int, postId int) (models.FollowedPost, error)
	DeleteFollowedPost(followedPostId int) error
	GetAllFollowedPost(userId int) ([]models.FollowedPost, error)
}
