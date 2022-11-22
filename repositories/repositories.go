package repositories

import (
	"discusiin/models"
)

type IDatabase interface {
	SaveNewUser(models.User) error
	GetUserByEmail(email string) (models.User, error)
	GetUserByUsername(username string) (models.User, error)
}
