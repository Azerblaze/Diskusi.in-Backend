package comments

import (
	"discusiin/models"
)

type ICommentDatabase interface {
	SaveNewComment(comment models.Comment) error
	GetAllCommentByPost(postID int) ([]models.Comment, error)
	GetCommentById(commendID int) (models.Comment, error)
	GetCommentByUserId(userId int, page int) ([]models.Comment, error)
	SaveComment(comment models.Comment) error
	DeleteComment(commentID int) error
	CountCommentByUserID(userId int) (int, error)
}
