package repositories

import (
	"discusiin/models"
)

type IDatabase interface {
	SaveNewUser(models.User) error
	GetUserByUsername(username string) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	GetUserById(userId int) (models.User, error)
	GetUsers(page int) ([]models.User, error)
	GetProfile(id int) (models.User, error)
	UpdateProfile(user models.User) error
	DeleteUser(userId int) error

	GetAllTopics() ([]models.Topic, error)
	GetTopicByName(topicName string) (models.Topic, error)
	GetTopicByID(topicID int) (models.Topic, error)
	SaveNewTopic(models.Topic) error
	SaveTopic(models.Topic) error
	RemoveTopic(topicID int) error

	SaveNewPost(post models.Post) error
	GetAllPostByTopic(topidID int, page int, search string) ([]models.Post, error)
	GetPostById(postID int) (models.Post, error)
	GetRecentPost(page int, search string) ([]models.Post, error)
	GetPostByUserId(userId int, page int) ([]models.Post, error)
	GetAllPostByLike(page int) ([]models.Post, error)
	SavePost(post models.Post) error
	DeletePost(postID int) error
	GetPostByIdWithAll(postID int) (models.Post, error)

	SaveNewComment(comment models.Comment) error
	GetAllCommentByPost(postID int) ([]models.Comment, error)
	GetCommentById(commendID int) (models.Comment, error)
	GetCommentByUserId(userId int, page int) ([]models.Comment, error)
	SaveComment(comment models.Comment) error
	DeleteComment(commentID int) error

	SaveNewReply(reply models.Reply) error
	GetAllReplyByComment(commentId int) ([]models.Reply, error)
	GetReplyById(re int) (models.Reply, error)
	SaveReply(reply models.Reply) error
	DeleteReply(re int) error

	GetLikeByUserAndPostId(userId int, postId int) (models.Like, error)
	SaveNewLike(like models.Like) error
	SaveLike(like models.Like) error

	SaveFollowedPost(followedPost models.FollowedPost) error
	GetFollowedPost(userId int, postId int) (models.FollowedPost, error)
	DeleteFollowedPost(followedPostId int) error
	GetAllFollowedPost(userId int) ([]models.FollowedPost, error)

	SaveBookmark(bookmark models.Bookmark) error
	GetBookmarkByUserIDAndPostID(userID, postID int) (models.Bookmark, error)
	GetBookmarkByBookmarkID(bookmarkID int) (models.Bookmark, error)
	DeleteBookmark(bookmarkId int) error
	GetAllBookmark(userId int) ([]models.Bookmark, error)

	CountPostLike(postID int) (int, error)
	CountPostComment(postID int) (int, error)
	CountPostDislike(postID int) (int, error)
	CountAllPost() (int, error)
	CountPostByTopicID(topicId int) (int, error)
	CountPostByUserID(userId int) (int, error)
	CountCommentByUserID(userId int) (int, error)
}
