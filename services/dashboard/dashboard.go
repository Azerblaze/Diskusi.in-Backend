package dashboard

import (
	"discusiin/dto"
	"discusiin/repositories"
	"net/http"

	"github.com/labstack/echo/v4"
)

func NewDashboardServices(db repositories.IDatabase) IDashboardServices {
	return &dashboardServices{IDatabase: db}
}

type IDashboardServices interface {
	GetAllTotal(token dto.Token) (int, int, int, error)
}

type dashboardServices struct {
	repositories.IDatabase
}

func (d *dashboardServices) GetAllTotal(token dto.Token) (int, int, int, error) {
	//check user
	user, errGetUser := d.IDatabase.GetUserByUsername(token.Username)
	if errGetUser != nil {
		if errGetUser.Error() == "record not found" {
			return 0, 0, 0, echo.NewHTTPError(http.StatusNotFound, "User not found")
		} else {
			return 0, 0, 0, echo.NewHTTPError(http.StatusInternalServerError, errGetUser.Error())
		}
	}
	if !user.IsAdmin {
		return 0, 0, 0, echo.NewHTTPError(http.StatusForbidden, "Admin access only")
	}

	//get user total
	userCount, errUserCount := d.IDatabase.CountAllUser()
	if errUserCount != nil {
		return 0, 0, 0, echo.NewHTTPError(http.StatusInternalServerError, errUserCount.Error())
	}
	//get topic total
	topicCount, errTopicCount := d.IDatabase.CountAllTopic()
	if errTopicCount != nil {
		return 0, 0, 0, echo.NewHTTPError(http.StatusInternalServerError, errTopicCount.Error())
	}
	//get post total
	postCount, errPostCount := d.IDatabase.CountAllPost()
	if errPostCount != nil {
		return 0, 0, 0, echo.NewHTTPError(http.StatusInternalServerError, errPostCount.Error())
	}

	return userCount, topicCount, postCount, nil
}
