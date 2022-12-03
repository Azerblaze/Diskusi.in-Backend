package users

import (
	"discusiin/dto"
	"discusiin/helper"
	"discusiin/middleware"
	"discusiin/models"
	"discusiin/repositories"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func NewUserServices(db repositories.IDatabase) IUserServices {
	return &userServices{IDatabase: db}
}

type IUserServices interface {
	Register(user models.User) error
	Login(user models.User) (dto.Login, error)
	GetUsers(token dto.Token, page int) ([]dto.PublicUser, error)
}

type userServices struct {
	repositories.IDatabase
}

func (s *userServices) Register(user models.User) error {
	var (
		client        models.User
		usernameTaken = true
		emailTaken    = true
	)

	client.Username = strings.ToLower(user.Username)
	_, errCheckUsername := s.IDatabase.GetUserByUsername(client.Username)
	if errCheckUsername != nil {
		if errCheckUsername.Error() == "record not found" {
			usernameTaken = false
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, errCheckUsername.Error())
		}
	} else {
		return echo.NewHTTPError(http.StatusConflict, "username has been taken")
	}
	if !usernameTaken {
		client.Email = strings.ToLower(user.Email)
		_, errCheckEmail := s.IDatabase.GetUserByEmail(client.Email)
		if errCheckEmail != nil {
			if errCheckEmail.Error() == "record not found" {
				emailTaken = false
			} else {
				return echo.NewHTTPError(http.StatusInternalServerError, errCheckEmail.Error())
			}
		} else {
			return echo.NewHTTPError(http.StatusConflict, "email has been used in another account")
		}
	}
	if !emailTaken {
		hashedPWD, errHashPassword := helper.HashPassword(user.Password)
		if errHashPassword != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, errHashPassword.Error())
		}
		client.Password = hashedPWD
		client.IsAdmin = user.IsAdmin
		err := s.IDatabase.SaveNewUser(client)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return nil
}
func (s *userServices) Login(user models.User) (dto.Login, error) {

	data, err := s.IDatabase.GetUserByEmail(user.Email)
	if err != nil {
		if err.Error() == "record not found" {
			return dto.Login{}, echo.NewHTTPError(http.StatusNotFound, "no account using this email")
		}
		return dto.Login{}, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var result dto.Login
	valid := helper.CheckPasswordHash(user.Password, data.Password)
	if valid {
		token, err := middleware.GetToken(data.ID, data.Username)
		if err != nil {
			return dto.Login{}, echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		result = dto.Login{
			ID:       data.ID,
			Username: data.Username,
			Email:    data.Email,
			Photo:    data.Photo,
			BanUntil: data.BanUntil,
			IsAdmin:  data.IsAdmin,
			Token:    token,
		}
	} else {
		return dto.Login{}, echo.NewHTTPError(http.StatusForbidden, "password incorrect")
	}
	return result, nil
}
func (s *userServices) GetUsers(token dto.Token, page int) ([]dto.PublicUser, error) {
	u, errGetUserByUsername := s.IDatabase.GetUserByUsername(token.Username)
	if errGetUserByUsername != nil {
		if errGetUserByUsername.Error() == "record not found" {
			return nil, echo.NewHTTPError(http.StatusNotFound, "your JWT does not have username field")
		} else {
			return nil, echo.NewHTTPError(http.StatusInternalServerError, errGetUserByUsername.Error())
		}
	}

	if !u.IsAdmin {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "admin access only")
	}
	users, err := s.IDatabase.GetUsers(page)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var result []dto.PublicUser
	for _, user := range users {
		result = append(result, dto.PublicUser{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Photo:    user.Photo,
			BanUntil: user.BanUntil,
			IsAdmin:  user.IsAdmin,
		})
	}
	return result, nil
}
