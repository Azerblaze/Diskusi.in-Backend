package users

import (
	"discusiin/dto"
	"discusiin/helper"
	"discusiin/middleware"
	"discusiin/models"
	"discusiin/repositories"
	"errors"
)

const RECORD_NOT_FOUND = "record not found"

func NewUserServices(db repositories.IDatabase) IUserServices {
	return &userServices{IDatabase: db}
}

type IUserServices interface {
	Register(user models.User) error
	Login(user models.User) (dto.Login, error)
}

type userServices struct {
	repositories.IDatabase
}

func (s *userServices) Register(user models.User) error {

	// isEmailTaken?
	_, errorCheckEmail := s.IDatabase.GetUserByEmail(user.Email)
	if errorCheckEmail != nil {
		if errorCheckEmail.Error() != RECORD_NOT_FOUND {
			return errorCheckEmail
		}
	} else {
		return errors.New("email has been taken")
	}
	// isUsernameTaken?
	_, errorCheckUname := s.IDatabase.GetUserByUsername(user.Username)
	if errorCheckUname != nil {
		if errorCheckUname.Error() == RECORD_NOT_FOUND {
			// bcrypt password | hash password
			hashedPWD, errorHash := helper.HashPassword(user.Password)
			user.Password = hashedPWD
			if errorHash != nil {
				return errorHash
			}
			// set !isAdmin
			user.IsAdmin = false
			// save to db
			errorSave := s.IDatabase.SaveNewUser(user)
			if errorSave != nil {
				return errorSave
			}
		} else {
			return errorCheckUname
		}
	} else {
		// if getUserByUsername got no error
		return errors.New("username has been taken")
	}
	return nil
}

func (s *userServices) Login(u models.User) (dto.Login, error) {

	user, err := s.IDatabase.GetUserByEmail(u.Email)
	if err != nil {
		return dto.Login{}, err
	}

	valid := helper.CheckPasswordHash(u.Password, user.Password)
	if valid {
		token, err := middleware.GetToken(user.ID, user.Username)
		if err != nil {
			return dto.Login{}, err
		}

		var result dto.Login
		result.ID = user.ID
		result.Username = user.Username
		result.Token = token

		return result, nil
	}
	return dto.Login{}, errors.New("password did not match")
}
