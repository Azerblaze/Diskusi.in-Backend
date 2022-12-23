package dashboard

import (
	"discusiin/dto"
	"discusiin/repositories"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func NewDashboardServices(userRepo repositories.IUserRepository, topicRepo repositories.ITopicRepository, postRepo repositories.IPostRepository) IDashboardServices {
	return &dashboardServices{IUserRepository: userRepo, ITopicRepository: topicRepo, IPostRepository: postRepo}
}

type IDashboardServices interface {
	GetTotalCountOfUserAndTopicAndPost(token dto.Token) (int, int, int, error)
}

type dashboardServices struct {
	repositories.IUserRepository
	repositories.ITopicRepository
	repositories.IPostRepository
}

func (d *dashboardServices) GetTotalCountOfUserAndTopicAndPost(token dto.Token) (int, int, int, error) {
	//check user
	user, errGetUser := d.IUserRepository.GetUserByUsername(token.Username)
	if errors.Is(errGetUser, gorm.ErrRecordNotFound) {
		return 0, 0, 0, echo.NewHTTPError(http.StatusNotFound, "User not found")
	} else if errGetUser != nil {
		return 0, 0, 0, echo.NewHTTPError(http.StatusInternalServerError, errGetUser.Error())
	}

	if !user.IsAdmin {
		return 0, 0, 0, echo.NewHTTPError(http.StatusForbidden, "Admin access only")
	}

	//get user total
	userCount, errUserCount := d.IUserRepository.CountAllUserNotAdminNotIncludeDeletedUser()
	if errUserCount != nil {
		return 0, 0, 0, echo.NewHTTPError(http.StatusInternalServerError, errUserCount.Error())
	}
	//get topic total
	topicCount, errTopicCount := d.ITopicRepository.CountAllTopic()
	if errTopicCount != nil {
		return 0, 0, 0, echo.NewHTTPError(http.StatusInternalServerError, errTopicCount.Error())
	}
	//get post total
	postCount, errPostCount := d.IPostRepository.CountAllPost()
	if errPostCount != nil {
		return 0, 0, 0, echo.NewHTTPError(http.StatusInternalServerError, errPostCount.Error())
	}

	return userCount, topicCount, postCount, nil
}
