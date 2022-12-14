package users

import (
	"discusiin/dto"
	"discusiin/helper"
	"discusiin/middleware"
	"discusiin/models"
	"discusiin/repositories"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func NewUserServices(db repositories.IDatabase) IUserServices {
	return &userServices{IDatabase: db}
}

type IUserServices interface {
	Register(user models.User) error
	RegisterAdmin(user models.User, token dto.Token) error
	Login(user models.User) (dto.Login, error)
	GetUsers(token dto.Token, page int) ([]dto.PublicUser, error)
	GetProfile(token dto.Token, user models.User) (dto.PublicUser, error)
	UpdateProfile(token dto.Token, user models.User) error
	DeleteUser(token dto.Token, userId int) error
	GetPostAsAdmin(token dto.Token, userId int, page int) (models.User, []dto.PublicPost, int, error)
	GetPostAsUser(token dto.Token, page int) ([]dto.PublicPost, int, error)
	BanUser(token dto.Token, userId int, user models.User) error
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
		return echo.NewHTTPError(http.StatusConflict, "Username has been taken")
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
			return echo.NewHTTPError(http.StatusConflict, "Email has been used in another account")
		}
	}
	//check if user registered as admin
	if user.IsAdmin {
		return echo.NewHTTPError(http.StatusForbidden, "Admin access only")
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

func (s *userServices) RegisterAdmin(user models.User, token dto.Token) error {
	//check user
	userAdmin, errGetUser := s.IDatabase.GetUserByUsername(token.Username)
	if errGetUser != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errGetUser.Error())
	}

	//check if user are admin
	if !userAdmin.IsAdmin {
		return echo.NewHTTPError(http.StatusForbidden, "Admin access only")
	}

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
		return echo.NewHTTPError(http.StatusConflict, "Username has been taken")
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
			return echo.NewHTTPError(http.StatusConflict, "Email has been used in another account")
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
			return dto.Login{}, echo.NewHTTPError(http.StatusNotFound, "Email or Password incorrect")
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
		return dto.Login{}, echo.NewHTTPError(http.StatusForbidden, "Email or Password incorrect")
	}

	var ban int
	//check if user are not banned
	if data.BanUntil > int(time.Now().UnixMilli()) {
		banLeft := data.BanUntil - int(time.Now().UnixMilli())
		ban = banLeft / 86400000

		//jika ban kurang dari 24 jam
		if ban < 1 {
			ban = banLeft / 3600
			return dto.Login{}, echo.NewHTTPError(http.StatusForbidden, "Ban Left: "+strconv.Itoa(ban)+" Hours")
		}

		return dto.Login{}, echo.NewHTTPError(http.StatusForbidden, "Ban Left: "+strconv.Itoa(ban)+" Days")
	}

	return result, nil
}
func (s *userServices) GetUsers(token dto.Token, page int) ([]dto.PublicUser, error) {
	u, errGetUserByUsername := s.IDatabase.GetUserByUsername(token.Username)
	if errGetUserByUsername != nil {
		if errGetUserByUsername.Error() == "record not found" {
			return nil, echo.NewHTTPError(http.StatusNotFound, "Invalid JWT Data")
		} else {
			return nil, echo.NewHTTPError(http.StatusInternalServerError, errGetUserByUsername.Error())
		}
	}

	if !u.IsAdmin {
		return nil, echo.NewHTTPError(http.StatusForbidden, "Admin access only")
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

func (s *userServices) GetProfile(token dto.Token, u models.User) (dto.PublicUser, error) {
	user, errGetProfile := s.IDatabase.GetProfile(int(token.ID))
	if errGetProfile != nil {
		if errGetProfile.Error() == "record not found" {
			return dto.PublicUser{}, echo.NewHTTPError(http.StatusNotFound, "Invalid JWT Data")
		} else {
			return dto.PublicUser{}, echo.NewHTTPError(http.StatusInternalServerError, errGetProfile.Error())
		}
	}
	result := dto.PublicUser{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Photo:    user.Photo,
		IsAdmin:  user.IsAdmin,
		BanUntil: user.BanUntil,
	}

	return result, nil
}

func (s *userServices) UpdateProfile(token dto.Token, user models.User) error {
	//get old profile
	oldProfile, errGetProfile := s.IDatabase.GetProfile(int(token.ID))
	if errGetProfile != nil {
		if errGetProfile.Error() == "record not found" {
			return echo.NewHTTPError(http.StatusNotFound, "Invalid JWT Data")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, errGetProfile.Error())
		}
	}

	oldProfile.Username = user.Username
	oldProfile.Photo = user.Photo

	//update profile
	errUpdateProfile := s.IDatabase.UpdateProfile(oldProfile)
	if errUpdateProfile != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errUpdateProfile.Error())
	}

	return nil
}

func (s *userServices) DeleteUser(token dto.Token, userId int) error {
	//check user
	user, err := s.IDatabase.GetUserByUsername(token.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if !user.IsAdmin {
		return echo.NewHTTPError(http.StatusForbidden, "Admin access only")
	}

	errDeleteUser := s.IDatabase.DeleteUser(userId)
	if errDeleteUser != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errDeleteUser.Error())
	}

	return nil
}

func (s *userServices) GetPostAsAdmin(token dto.Token, userId int, page int) (models.User, []dto.PublicPost, int, error) {
	//check user Admin
	userAdmin, errUserAdmin := s.IDatabase.GetUserByUsername(token.Username)
	if errUserAdmin != nil {
		return models.User{}, nil, 0, echo.NewHTTPError(http.StatusInternalServerError, errUserAdmin.Error())
	}

	//check if logged user is admin
	if !userAdmin.IsAdmin {
		return models.User{}, nil, 0, echo.NewHTTPError(http.StatusForbidden, "Admin access only")
	}

	//check user
	user, errUser := s.IDatabase.GetUserById(userId)
	if errUser != nil {
		if errUser.Error() == "record not found" {
			return models.User{}, nil, 0, echo.NewHTTPError(http.StatusNotFound, "User not found")
		} else {
			return models.User{}, nil, 0, echo.NewHTTPError(http.StatusInternalServerError, errUser.Error())
		}
	}
	user.Password = "<secret>"

	//cek jika page kosong
	if page < 1 {
		page = 1
	}

	//get post by user id
	posts, errGetPostByUserId := s.IDatabase.GetPostByUserId(userId, page)
	if errGetPostByUserId != nil {
		return models.User{}, nil, 0, echo.NewHTTPError(http.StatusInternalServerError, errGetPostByUserId.Error())
	}

	//insert data to dto.Public post
	var result []dto.PublicPost
	for _, post := range posts {
		likeCount, _ := s.IDatabase.CountPostLike(int(post.ID))
		commentCount, _ := s.IDatabase.CountPostComment(int(post.ID))
		dislikeCount, _ := s.IDatabase.CountPostDislike(int(post.ID))

		result = append(result, dto.PublicPost{
			Model:     post.Model,
			Title:     post.Title,
			Photo:     post.Photo,
			Body:      post.Body,
			CreatedAt: post.CreatedAt,
			IsActive:  post.IsActive,
			User: dto.PostUser{
				UserID:   post.UserID,
				Photo:    post.User.Photo,
				Username: post.User.Username,
			},
			Topic: dto.PostTopic{
				TopicID:   post.TopicID,
				TopicName: post.Topic.Name,
			},
			Count: dto.PostCount{
				LikeCount:    likeCount,
				CommentCount: commentCount,
				DislikeCount: dislikeCount,
			},
		})
	}

	//count page number
	numberOfPost, errPage := s.IDatabase.CountPostByUserID(userId)
	if errPage != nil {
		return models.User{}, nil, 0, echo.NewHTTPError(http.StatusInternalServerError, errPage.Error())
	}

	//count number of page
	var numberOfPage int
	if numberOfPost%20 == 0 {
		numberOfPage = (numberOfPost / 20)
	} else {
		numberOfPage = (numberOfPost / 20) + 1
	}

	return user, result, numberOfPage, nil
}

func (s *userServices) GetPostAsUser(token dto.Token, page int) ([]dto.PublicPost, int, error) {
	//cek jika page kosong
	if page < 1 {
		page = 1
	}

	//get post by user id
	posts, errGetPostByUserId := s.IDatabase.GetPostByUserId(int(token.ID), page)
	if errGetPostByUserId != nil {
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError, errGetPostByUserId.Error())
	}

	//insert data to dto.Public post
	var result []dto.PublicPost
	for _, post := range posts {
		likeCount, _ := s.IDatabase.CountPostLike(int(post.ID))
		commentCount, _ := s.IDatabase.CountPostComment(int(post.ID))
		dislikeCount, _ := s.IDatabase.CountPostDislike(int(post.ID))

		result = append(result, dto.PublicPost{
			Model:     post.Model,
			Title:     post.Title,
			Photo:     post.Photo,
			Body:      post.Body,
			CreatedAt: post.CreatedAt,
			IsActive:  post.IsActive,
			User: dto.PostUser{
				UserID:   post.UserID,
				Photo:    post.User.Photo,
				Username: post.User.Username,
			},
			Topic: dto.PostTopic{
				TopicID:   post.TopicID,
				TopicName: post.Topic.Name,
			},
			Count: dto.PostCount{
				LikeCount:    likeCount,
				CommentCount: commentCount,
				DislikeCount: dislikeCount,
			},
		})
	}

	//count page number
	numberOfPost, errPage := s.IDatabase.CountPostByUserID(int(token.ID))
	if errPage != nil {
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError, errPage.Error())
	}

	//count number of page
	var numberOfPage int
	if numberOfPost%20 == 0 {
		numberOfPage = (numberOfPost / 20)
	} else {
		numberOfPage = (numberOfPost / 20) + 1
	}

	return result, numberOfPage, nil
}

func (s *userServices) BanUser(token dto.Token, userId int, user models.User) error {
	//check user admin
	userAdmin, errUserAdmin := s.IDatabase.GetUserByUsername(token.Username)
	if errUserAdmin != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errUserAdmin.Error())
	}

	//check if logged user is admin
	if !userAdmin.IsAdmin {
		return echo.NewHTTPError(http.StatusForbidden, "Admin access only")
	}

	//check if user exist
	oldUser, errUser := s.IDatabase.GetUserById(userId)
	if errUser != nil {
		if errUser.Error() == "record not found" {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, errUser.Error())
		}
	}

	//user variabel ban to store how long ban wil last
	ban := user.BanUntil
	const DAY_IN_UNIX_MILLISECOND = 86400000
	user.BanUntil = int(time.Now().UnixMilli()) + (DAY_IN_UNIX_MILLISECOND * ban)

	//update user
	oldUser.BanUntil = user.BanUntil
	err := s.IDatabase.UpdateProfile(oldUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}
