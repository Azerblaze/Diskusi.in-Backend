package replies

import (
	"discusiin/models"
)

type IReplyDatabase interface {
	SaveNewReply(reply models.Reply) error
	GetAllReplyByComment(commentId int) ([]models.Reply, error)
	GetReplyById(re int) (models.Reply, error)
	SaveReply(reply models.Reply) error
	DeleteReply(re int) error
}
