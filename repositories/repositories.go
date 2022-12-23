package repositories

import (
	"discusiin/dto"
	"discusiin/models"
)

type IUserRepository interface {
	SaveNewUser(models.User) error
	GetUserByUsername(username string) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	GetUserById(userId int) (models.User, error)
	GetUsersAdminNotIncluded(page int) ([]models.User, error)
	GetProfile(id int) (models.User, error)
	UpdateProfile(user models.User) error
	DeleteUser(userId int) error
	CountAllUserNotIncludeDeletedUser() (int, error)
	CountAllUserNotAdminNotIncludeDeletedUser() (int, error)
}
type IPostRepository interface {
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
}
type ITopicRepository interface {
	GetAllTopics() ([]models.Topic, error)
	GetTopicByName(topicName string) (models.Topic, error)
	GetTopicByID(topicID int) (models.Topic, error)
	GetTopTopics() ([]dto.TopTopics, error)
	SaveNewTopic(models.Topic) error
	SaveTopic(models.Topic) error
	RemoveTopic(topicID int) error
	CountAllTopic() (int, error)
	CountNumberOfPostByTopicName(topicName string) (int, error)
}
type ICommentRepository interface {
	SaveNewComment(comment models.Comment) error
	GetAllCommentByPost(postID int) ([]models.Comment, error)
	GetCommentById(commendID int) (models.Comment, error)
	GetCommentByUserId(userId int, page int) ([]models.Comment, error)
	SaveComment(comment models.Comment) error
	DeleteComment(commentID int) error
	CountCommentByUserID(userId int) (int, error)
}
type IReplyRepository interface {
	SaveNewReply(reply models.Reply) error
	GetAllReplyByComment(commentId int) ([]models.Reply, error)
	GetReplyById(re int) (models.Reply, error)
	SaveReply(reply models.Reply) error
	DeleteReply(re int) error
}
type ILikeRepository interface {
	GetLikeByUserAndPostId(userId int, postId int) (models.Like, error)
	SaveNewLike(like models.Like) error
	SaveLike(like models.Like) error
}
type IFollowedPostRepository interface {
	SaveFollowedPost(followedPost models.FollowedPost) error
	GetFollowedPost(userId int, postId int) (models.FollowedPost, error)
	DeleteFollowedPost(followedPostId int) error
	GetAllFollowedPost(userId int) ([]models.FollowedPost, error)
}
type IBookmarkRepository interface {
	SaveBookmark(bookmark models.Bookmark) error
	GetBookmarkByUserIDAndPostID(userID, postID int) (models.Bookmark, error)
	GetBookmarkByBookmarkID(bookmarkID int) (models.Bookmark, error)
	DeleteBookmark(bookmarkId int) error
	GetAllBookmark(userId int) ([]models.Bookmark, error)
}
