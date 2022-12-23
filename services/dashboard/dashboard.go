package dashboard

import (
	"discusiin/dto"
	"discusiin/repositories/posts"
	"discusiin/repositories/topics"
	"discusiin/repositories/users"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func NewDashboardServices(dbUser users.IUserDatabase, dbPost posts.IPostDatabase, dbTopic topics.ITopicDatabase) IDashboardServices {
	return &dashboardServices{IUserDatabase: dbUser, IPostDatabase: dbPost, ITopicDatabase: dbTopic}
}

type IDashboardServices interface {
	GetTotalCountOfUserAndTopicAndPost(token dto.Token) (int, int, int, error)
}

type dashboardServices struct {
	users.IUserDatabase
	posts.IPostDatabase
	topics.ITopicDatabase
}

func (d *dashboardServices) GetTotalCountOfUserAndTopicAndPost(token dto.Token) (int, int, int, error) {
	//check user
	user, errGetUser := d.IUserDatabase.GetUserByUsername(token.Username)
	if errGetUser != nil {
		if errGetUser == gorm.ErrRecordNotFound {
			return 0, 0, 0, echo.NewHTTPError(http.StatusNotFound, "User not found")
		} else {
			return 0, 0, 0, echo.NewHTTPError(http.StatusInternalServerError, errGetUser.Error())
		}
	}
	if !user.IsAdmin {
		return 0, 0, 0, echo.NewHTTPError(http.StatusForbidden, "Admin access only")
	}

	//get user total
	userCount, errUserCount := d.IUserDatabase.CountAllUserNotAdminNotIncludeDeletedUser()
	if errUserCount != nil {
		return 0, 0, 0, echo.NewHTTPError(http.StatusInternalServerError, errUserCount.Error())
	}
	//get topic total
	topicCount, errTopicCount := d.ITopicDatabase.CountAllTopic()
	if errTopicCount != nil {
		return 0, 0, 0, echo.NewHTTPError(http.StatusInternalServerError, errTopicCount.Error())
	}
	//get post total
	postCount, errPostCount := d.IPostDatabase.CountAllPost()
	if errPostCount != nil {
		return 0, 0, 0, echo.NewHTTPError(http.StatusInternalServerError, errPostCount.Error())
	}

	return userCount, topicCount, postCount, nil
}
