package users

import (
	"discusiin/dto"
	"discusiin/helper"
	"discusiin/middleware"
	"discusiin/models"
	"discusiin/repositories/comments"
	"discusiin/repositories/posts"
	"discusiin/repositories/users"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func NewUserServices(dbUser users.IUserDatabase, dbPost posts.IPostDatabase, dbComment comments.ICommentDatabase) IUserServices {
	return &userServices{IUserDatabase: dbUser, IPostDatabase: dbPost, ICommentDatabase: dbComment}
}

type IUserServices interface {
	Register(user models.User) error
	RegisterAdmin(user models.User, token dto.Token) error
	Login(user models.User) (dto.Login, error)
	GetUsersAdminNotIncluded(token dto.Token, page int) ([]dto.PublicUser, int, error)
	GetProfile(token dto.Token, user models.User) (dto.PublicUser, error)
	UpdateProfile(token dto.Token, user models.User) error
	DeleteUser(token dto.Token, userId int) error
	GetPostAsAdmin(token dto.Token, userId int, page int) (models.User, []dto.PublicPost, int, error)
	GetCommentAsAdmin(token dto.Token, userId int, page int) (models.User, []dto.AdminComment, int, error)
	GetPostAsUser(token dto.Token, page int) ([]dto.PublicPost, int, error)
	BanUser(token dto.Token, userId int, user models.User) (dto.PublicUser, error)
}

type userServices struct {
	users.IUserDatabase
	posts.IPostDatabase
	comments.ICommentDatabase
}

func (s *userServices) Register(user models.User) error {
	var client models.User

	//check if user registered as admin
	if user.IsAdmin {
		return echo.NewHTTPError(http.StatusForbidden, "Admin access only")
	}
	client.Username = strings.ToLower(user.Username)
	client.Email = strings.ToLower(user.Email)
	hashedPWD, errHashPassword := helper.HashPassword(user.Password)
	if errHashPassword != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errHashPassword.Error())
	}
	client.Password = hashedPWD
	client.IsAdmin = user.IsAdmin

	_, errCheckUsername := s.IUserDatabase.GetUserByUsername(client.Username)
	if errCheckUsername != nil {
		if errCheckUsername == gorm.ErrRecordNotFound {
			_, errCheckEmail := s.IUserDatabase.GetUserByEmail(client.Email)
			if errCheckEmail != nil {
				if errCheckEmail == gorm.ErrRecordNotFound {
					err := s.IUserDatabase.SaveNewUser(client)
					if err != nil {
						return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
					}
				} else {
					return echo.NewHTTPError(http.StatusInternalServerError, errCheckEmail.Error())
				}
			} else {
				return echo.NewHTTPError(http.StatusConflict, "Email has been used in another account")
			}
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, errCheckUsername.Error())
		}
	} else {
		return echo.NewHTTPError(http.StatusConflict, "Username has been taken")
	}
	return nil
}

func (s *userServices) RegisterAdmin(user models.User, token dto.Token) error {
	//check user
	userAdmin, errGetUser := s.IUserDatabase.GetUserByUsername(token.Username)
	if errGetUser != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errGetUser.Error())
	}

	//check if user are admin
	if !userAdmin.IsAdmin {
		return echo.NewHTTPError(http.StatusForbidden, "Admin access only")
	}

	var (
		client models.User
	)

	client.Username = strings.ToLower(user.Username)
	client.Email = strings.ToLower(user.Email)
	hashedPWD, errHashPassword := helper.HashPassword(user.Password)
	if errHashPassword != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errHashPassword.Error())
	}
	client.Password = hashedPWD
	client.IsAdmin = user.IsAdmin

	_, errCheckUsername := s.IUserDatabase.GetUserByUsername(client.Username)
	if errCheckUsername != nil {
		if errCheckUsername == gorm.ErrRecordNotFound {
			_, errCheckEmail := s.IUserDatabase.GetUserByEmail(client.Email)
			if errCheckEmail != nil {
				if errCheckEmail == gorm.ErrRecordNotFound {
					err := s.IUserDatabase.SaveNewUser(client)
					if err != nil {
						return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
					}
				} else {
					return echo.NewHTTPError(http.StatusInternalServerError, errCheckEmail.Error())
				}
			} else {
				return echo.NewHTTPError(http.StatusConflict, "Email has been used in another account")
			}
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, errCheckUsername.Error())
		}
	} else {
		return echo.NewHTTPError(http.StatusConflict, "Username has been taken")
	}
	return nil
}
func (s *userServices) Login(user models.User) (dto.Login, error) {

	data, err := s.IUserDatabase.GetUserByEmail(user.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
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
func (s *userServices) GetUsersAdminNotIncluded(token dto.Token, page int) ([]dto.PublicUser, int, error) {
	u, errGetUserByUsername := s.IUserDatabase.GetUserByUsername(token.Username)
	if errGetUserByUsername != nil {
		if errGetUserByUsername == gorm.ErrRecordNotFound {
			return nil, 0, echo.NewHTTPError(http.StatusNotFound, "Invalid JWT Data")
		} else {
			return nil, 0, echo.NewHTTPError(http.StatusInternalServerError, errGetUserByUsername.Error())
		}
	}

	if !u.IsAdmin {
		return nil, 0, echo.NewHTTPError(http.StatusForbidden, "Admin access only")
	}
	users, err := s.IUserDatabase.GetUsersAdminNotIncluded(page)
	if err != nil {
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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
	userCount, _ := s.IUserDatabase.CountAllUserNotIncludeDeletedUser()
	// Jumlah data per page
	pageSize := 20

	// Hitung jumlah page dengan pembagian sederhana
	numberOfPage := math.Ceil(float64(userCount) / float64(pageSize))

	// Jika ada sisa, tambahkan 1 page untuk menampung sisa data tersebut

	return result, int(numberOfPage), nil
}

func (s *userServices) GetProfile(token dto.Token, u models.User) (dto.PublicUser, error) {
	user, errGetProfile := s.IUserDatabase.GetProfile(int(token.ID))
	if errGetProfile != nil {
		if errGetProfile == gorm.ErrRecordNotFound {
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
	oldProfile, errGetProfile := s.IUserDatabase.GetProfile(int(token.ID))
	if errGetProfile != nil {
		if errGetProfile == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "Invalid JWT Data")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, errGetProfile.Error())
		}
	}

	oldProfile.Username = user.Username
	oldProfile.Photo = user.Photo

	//update profile
	errUpdateProfile := s.IUserDatabase.UpdateProfile(oldProfile)
	if errUpdateProfile != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errUpdateProfile.Error())
	}

	return nil
}

func (s *userServices) DeleteUser(token dto.Token, userId int) error {
	//check user admin
	userAdmin, err := s.IUserDatabase.GetUserByUsername(token.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if !userAdmin.IsAdmin {
		return echo.NewHTTPError(http.StatusForbidden, "Admin access only")
	}

	_, err = s.IUserDatabase.GetUserById(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	errDeleteUser := s.IUserDatabase.DeleteUser(userId)
	if errDeleteUser != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errDeleteUser.Error())
	}
	errDeletePost := s.IPostDatabase.DeletePostByUserID(userId)
	if errDeletePost != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errDeleteUser.Error())
	}

	return nil
}

func (s *userServices) GetCommentAsAdmin(token dto.Token, userId int, page int) (models.User, []dto.AdminComment, int, error) {
	//check user Admin
	userAdmin, errUserAdmin := s.IUserDatabase.GetUserByUsername(token.Username)
	if errUserAdmin != nil {
		return models.User{}, nil, 0, echo.NewHTTPError(http.StatusInternalServerError, errUserAdmin.Error())
	}

	//check if logged user is admin
	if !userAdmin.IsAdmin {
		return models.User{}, nil, 0, echo.NewHTTPError(http.StatusForbidden, "Admin access only")
	}

	//check user
	user, errUser := s.IUserDatabase.GetUserById(userId)
	if errUser != nil {
		if errUser == gorm.ErrRecordNotFound {
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

	//get comment by user id
	comments, errGetCommentByUserId := s.ICommentDatabase.GetCommentByUserId(userId, page)
	if errGetCommentByUserId != nil {
		return models.User{}, nil, 0, echo.NewHTTPError(http.StatusInternalServerError, errGetCommentByUserId.Error())
	}

	//insert data to dto.Public Comment
	var result []dto.AdminComment
	for _, comment := range comments {
		result = append(result, dto.AdminComment{
			Model: comment.Model,
			Body:  comment.Body,
			Post: dto.CommentPost{
				PostID: comment.PostID,
				Title:  comment.Post.Title,
				Body:   comment.Post.Body,
			},
		})
	}

	//count page number
	numberOfPost, errPage := s.ICommentDatabase.CountCommentByUserID(userId)
	if errPage != nil {
		return models.User{}, nil, 0, echo.NewHTTPError(http.StatusInternalServerError, errPage.Error())
	}

	// Jumlah data per page
	pageSize := 20

	// Hitung jumlah page dengan pembagian sederhana
	numberOfPage := math.Ceil(float64(numberOfPost) / float64(pageSize))

	// Jika ada sisa, tambahkan 1 page untuk menampung sisa data tersebut

	return user, result, int(numberOfPage), nil
}

func (s *userServices) GetPostAsAdmin(token dto.Token, userId int, page int) (models.User, []dto.PublicPost, int, error) {
	//check user Admin
	userAdmin, errUserAdmin := s.IUserDatabase.GetUserByUsername(token.Username)
	if errUserAdmin != nil {
		return models.User{}, nil, 0, echo.NewHTTPError(http.StatusInternalServerError, errUserAdmin.Error())
	}

	//check if logged user is admin
	if !userAdmin.IsAdmin {
		return models.User{}, nil, 0, echo.NewHTTPError(http.StatusForbidden, "Admin access only")
	}

	//check user
	user, errUser := s.IUserDatabase.GetUserById(userId)
	if errUser != nil {
		if errUser == gorm.ErrRecordNotFound {
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
	posts, errGetPostByUserId := s.IPostDatabase.GetPostByUserId(userId, page)
	if errGetPostByUserId != nil {
		return models.User{}, nil, 0, echo.NewHTTPError(http.StatusInternalServerError, errGetPostByUserId.Error())
	}

	//insert data to dto.Public post
	var result []dto.PublicPost
	for _, post := range posts {
		likeCount, _ := s.IPostDatabase.CountPostLike(int(post.ID))
		commentCount, _ := s.IPostDatabase.CountPostComment(int(post.ID))
		dislikeCount, _ := s.IPostDatabase.CountPostDislike(int(post.ID))

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
	numberOfPost, errPage := s.IPostDatabase.CountPostByUserID(userId)
	if errPage != nil {
		return models.User{}, nil, 0, echo.NewHTTPError(http.StatusInternalServerError, errPage.Error())
	}

	// Jumlah data per page
	pageSize := 20

	// Hitung jumlah page dengan pembagian sederhana
	numberOfPage := math.Ceil(float64(numberOfPost) / float64(pageSize))

	// Jika ada sisa, tambahkan 1 page untuk menampung sisa data tersebut

	return user, result, int(numberOfPage), nil
}

func (s *userServices) GetPostAsUser(token dto.Token, page int) ([]dto.PublicPost, int, error) {
	//cek jika page kosong
	if page < 1 {
		page = 1
	}

	//get post by user id
	posts, errGetPostByUserId := s.IPostDatabase.GetPostByUserId(int(token.ID), page)
	if errGetPostByUserId != nil {
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError, errGetPostByUserId.Error())
	}

	//insert data to dto.Public post
	var result []dto.PublicPost
	for _, post := range posts {
		likeCount, _ := s.IPostDatabase.CountPostLike(int(post.ID))
		commentCount, _ := s.IPostDatabase.CountPostComment(int(post.ID))
		dislikeCount, _ := s.IPostDatabase.CountPostDislike(int(post.ID))

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
	numberOfPost, errPage := s.IPostDatabase.CountPostByUserID(int(token.ID))
	if errPage != nil {
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError, errPage.Error())
	}

	// Jumlah data per page
	pageSize := 20

	// Hitung jumlah page dengan pembagian sederhana
	numberOfPage := math.Ceil(float64(numberOfPost) / float64(pageSize))

	// Jika ada sisa, tambahkan 1 page untuk menampung sisa data tersebut

	return result, int(numberOfPage), nil
}

func (s *userServices) BanUser(token dto.Token, userId int, user models.User) (dto.PublicUser, error) {
	//check user admin
	userAdmin, errUserAdmin := s.IUserDatabase.GetUserByUsername(token.Username)
	if errUserAdmin != nil {
		return dto.PublicUser{}, echo.NewHTTPError(http.StatusInternalServerError, errUserAdmin.Error())
	}

	//check if logged user is admin
	if !userAdmin.IsAdmin {
		return dto.PublicUser{}, echo.NewHTTPError(http.StatusForbidden, "Admin access only")
	}

	//check if user exist
	oldUser, errUser := s.IUserDatabase.GetUserById(userId)
	if errUser != nil {
		if errUser == gorm.ErrRecordNotFound {
			return dto.PublicUser{}, echo.NewHTTPError(http.StatusNotFound, "User not found")
		} else {
			return dto.PublicUser{}, echo.NewHTTPError(http.StatusInternalServerError, errUser.Error())
		}
	}

	//user variabel ban to store how long ban wil last
	ban := user.BanUntil
	const DayInUnixMillisecond = 86400000
	user.BanUntil = int(time.Now().UnixMilli()) + (DayInUnixMillisecond * ban)

	//update user
	oldUser.BanUntil = user.BanUntil
	err := s.IUserDatabase.UpdateProfile(oldUser)
	if err != nil {
		return dto.PublicUser{}, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	result := dto.PublicUser{
		ID:       oldUser.ID,
		Username: oldUser.Username,
		Email:    oldUser.Email,
		Photo:    oldUser.Photo,
		IsAdmin:  oldUser.IsAdmin,
		BanUntil: oldUser.BanUntil,
	}

	return result, nil
}
