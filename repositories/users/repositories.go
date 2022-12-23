package users

import (
	"discusiin/models"
)

type IUserDatabase interface {
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
