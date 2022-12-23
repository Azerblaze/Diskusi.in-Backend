package likes

import (
	"discusiin/models"
)

type ILikeDatabase interface {
	GetLikeByUserAndPostId(userId int, postId int) (models.Like, error)
	SaveNewLike(like models.Like) error
	SaveLike(like models.Like) error
}
